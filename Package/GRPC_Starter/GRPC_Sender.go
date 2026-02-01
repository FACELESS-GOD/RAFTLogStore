package GRPC_Starter

import (
	"context"
	"log"
	"sync"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	"google.golang.org/grpc"
)

func (Grc *GRPCService) Sender() {
	for {
		var is_Mode_Leader bool = false
		Grc.Ut.Mu.Lock()
		if Grc.Ut.Mode == State.Leader {
			is_Mode_Leader = true
		} else {
			is_Mode_Leader = false
		}
		Grc.Ut.Mu.Unlock()

		if is_Mode_Leader == true {
			Grc.DiscoverService()
			wg := sync.WaitGroup{}
			for _, serverAddress := range Grc.Service.ServerList {
				go Grc.UpdateTableNex(serverAddress, &wg)
				wg.Add(1)
			}
			wg.Wait()
		}
	}

}
func (Grc *GRPCService) DiscoverService() {}

func (Grc *GRPCService) UpdateTableNex(Address string, Wg *sync.WaitGroup) {
	defer Wg.Done()

	var rpcServiceClient GRPCServicePackage.RPCServiceClient

	_, isExits := Grc.TableConnClient[Address]

	if isExits == true {

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

	payload := GRPCServicePackage.AddLogRequest{
		IsHeartbeat: true,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Heart-Beat Message
	res, err := rpcServiceClient.AppendRPC(ctx, &payload)

	defer cancelFunc()

	if err != nil {
		log.Println(err)
	}

	currLog := Grc.TableNex[Address]
	childLog := res.GetLogId()

	if currLog != int(childLog) {
		if currLog < int(childLog) {

			Grc.mu.Lock()
			Grc.TableNex[Address] = int(childLog)
			Grc.mu.Unlock()

			Grc.UpdateLog(ctx, rpcServiceClient, int(childLog))
		}
	}

}

func (Grc *GRPCService) UpdateLog(Ctx context.Context, RpcServiceClient GRPCServicePackage.RPCServiceClient, ChildLogId int) {
	req := GRPCServicePackage.AddLogRequest{
		IsHeartbeat: false,
	}
	for i := ChildLogId; i < int(Grc.Ut.LogId); i++ {

		log := GRPCServicePackage.LogStructure{
			LogId:  int32(i),
			TermId: Grc.Ut.Term,
			Text:   Grc.Mdl.Arr[i],
		}

		req.Log = append(req.Log, &log)
	}

	_, err := RpcServiceClient.AppendRPC(Ctx, &req)
	if err != nil {
		log.Println(err)
	}
}
