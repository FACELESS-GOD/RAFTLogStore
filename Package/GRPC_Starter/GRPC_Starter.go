package GRPC_Starter

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/ServiceRegistry"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"google.golang.org/grpc"
)

type GRPCServiceInterface interface {
	AddLog(Log.LogStuct) (bool, error)
}

type GRPCService struct {
	AddRequest      chan Log.LogStuct
	Response        chan bool
	Ut              Util.UtilStruct
	Mdl             Model.ModelStuct
	Service         ServiceRegistry.Service
	TableNex        map[string]int
	mu              sync.Mutex
	ChildCount      int
	TableConn       map[string]*grpc.ClientConn
	TableConnClient map[string]GRPCServicePackage.RPCServiceClient
}

func NewGRPCService(req chan Log.LogStuct, res chan bool, Ut Util.UtilStruct, Mdl Model.ModelStuct) (GRPCService, error) {
	grc := GRPCService{
		AddRequest: req,
		Response:   res,
		Ut:         Ut,
		Mdl:        Mdl,
	}
	return grc, nil
}

func (Grc *GRPCService) AddLog(log Log.LogStuct) (bool, error) {
	go Grc.RecurAddLog()
	Grc.AddRequest <- log
	res := <-Grc.Response
	return res, nil
}

func (Grc *GRPCService) Run_GRPC_Server(Util Util.UtilStruct, ShutDownSignalChannel chan int) {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	server := Server{}
	server.Ut = Grc.Ut
	server.Mdl = Grc.Mdl

	GRPCServicePackage.RegisterRPCServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (Grc *GRPCService) RecurAddLog() {

	for {
		newLog := <-Grc.Mdl.AddLogChan

		for _, Address := range Grc.Service.ServerList {

			var rpcServiceClient GRPCServicePackage.RPCServiceClient

			_, is_Exists := Grc.TableConnClient[Address]

			if is_Exists == true {

				rpcServiceClient = Grc.TableConnClient[Address]

			} else {

				var conn *grpc.ClientConn
				conn, err := grpc.NewClient(Address, grpc.WithMaxCallAttempts(3))

				if err != nil {
					log.Println(err)
					continue
				}

				rpcServiceClient = GRPCServicePackage.NewRPCServiceClient(conn)

				Grc.mu.Lock()
				Grc.TableConn[Address] = conn
				Grc.TableConnClient[Address] = rpcServiceClient
				Grc.mu.Unlock()

			}

			payload := GRPCServicePackage.AddLogRequest{
				IsHeartbeat: false,
			}
			grpcLog := GRPCServicePackage.LogStructure{
				Text:   newLog.Text,
				LogId:  Grc.Ut.LogId,
				TermId: Grc.Ut.Term,
			}

			payload.Log = append(payload.Log, &grpcLog)

			ctx, cancelFunc := context.WithCancel(context.Background())

			res, err := rpcServiceClient.AppendRPC(ctx, &payload)

			defer cancelFunc()

			if err != nil {
				log.Println(err)
			}

			if res.IsAnyError == true {
				log.Print(res.ErrorMessages)
			}
		}

	}

}

func (Grc *GRPCService) ElectionsCheck() {

	for {
		currTime := time.Now()
		if currTime.Sub(Grc.Ut.LastTouch) > Grc.Ut.ElectionTimeout {

			Grc.Ut.Mu.Lock()
			Grc.Ut.Mode = State.Candidate
			Grc.Ut.Is_Voted = true
			Grc.Ut.Term = Grc.Ut.Term + 1
			Grc.ChildCount = 0
			Grc.Ut.Mu.Unlock()

			is_elected, is_tie := Grc.BeginElection()

			Grc.Ut.LastTouch = time.Now()

			if is_elected == true {

				Grc.Ut.Mu.Lock()
				Grc.Ut.Mode = State.Leader
				Grc.Ut.Mu.Unlock()

				wg := sync.WaitGroup{}
				for _, serverAddress := range Grc.Service.ServerList {
					go Grc.UpdateTableNex(serverAddress, &wg)
					wg.Add(1)
				}
				wg.Wait()

			} else {

				Grc.Ut.Mu.Lock()
				Grc.Ut.Mode = State.Follower
				Grc.Ut.Mu.Unlock()

				if is_tie == true {
					Grc.Ut.LastTouch = time.Now().Add(time.Millisecond * time.Duration(rand.Int()))
				}
			}

		}
	}

}

func (Grc *GRPCService) BeginElection() (bool, bool) {
	Grc.DiscoverService()
	wg := sync.WaitGroup{}
	for _, serverAddress := range Grc.Service.ServerList {
		go Grc.Elector(serverAddress, &wg)
		wg.Add(1)
	}
	wg.Wait()
	numOfChild := len(Grc.TableNex)
	boundryCount := numOfChild / 2
	if Grc.ChildCount == boundryCount {
		return false, true
	} else if Grc.ChildCount > boundryCount {
		return true, false
	} else {
		return false, false
	}

}

func (Grc *GRPCService) Elector(Address string, Wg *sync.WaitGroup) {
	defer Wg.Done()
	var rpcServiceClient GRPCServicePackage.RPCServiceClient

	_, is_Exists := Grc.TableConnClient[Address]

	if is_Exists == true {

		rpcServiceClient = Grc.TableConnClient[Address]

	} else {

		var conn *grpc.ClientConn
		conn, err := grpc.NewClient(Address, grpc.WithMaxCallAttempts(3))

		if err != nil {
			log.Println(err)
			return
		}

		rpcServiceClient = GRPCServicePackage.NewRPCServiceClient(conn)

		Grc.mu.Lock()
		Grc.TableConn[Address] = conn
		Grc.TableConnClient[Address] = rpcServiceClient
		Grc.mu.Unlock()

	}


	payload := GRPCServicePackage.RequestLogRequest{
		TermId: Grc.Ut.Term,
		LogId:  Grc.Ut.LogId,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Heart-Beat Message
	res, err := rpcServiceClient.RequestVoteRPC(ctx, &payload)

	defer cancelFunc()

	if err != nil {
		log.Println(err)
	}

	if res.Vote == true {
		Grc.mu.Lock()
		Grc.TableNex[Address] = int(res.LogId)
		Grc.ChildCount = Grc.ChildCount + 1
		Grc.mu.Lock()
	}

}
