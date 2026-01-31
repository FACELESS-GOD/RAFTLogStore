package GRPC_Starter

import (
	"fmt"
	"log"
	"net"
	"time"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"google.golang.org/grpc"
)

type GRPCServiceInterface interface {
	AddLog(Log.LogStuct) (bool, error)
}

type GRPCService struct {
	AddRequest chan Log.LogStuct
	Response   chan bool
	Ut         Util.UtilStruct
	Mdl        Model.ModelStuct
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
		fmt.Println(newLog.Text)
	}

}


func (Grc *GRPCService) ElectionsCheck() {

	for {
		currTime := time.Now()
		if currTime.Sub(Grc.Ut.LastTouch) > Grc.Ut.ElectionTimeout {
			Grc.Ut.Mu.Lock()
			Grc.Ut.Mode = State.Leader
			Grc.Ut.Is_Voted = true			
			Grc.Ut.Mu.Unlock()
			
			Grc.BeginElection()
		} 
	}

}

func (Grc *GRPCService) BeginElection() {



	
}
