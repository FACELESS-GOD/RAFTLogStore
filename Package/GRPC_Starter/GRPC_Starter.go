package GRPC_Starter

import (
	"log"
	"net"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"google.golang.org/grpc"
)

type GRPCServiceInterface interface {
	AddLog(Log.LogStuct) (bool, error)
}

type GRPCService struct {
	AddRequest chan Log.LogStuct
	Response   chan bool
	Ut Util.UtilStruct
}

func NewGRPCService(req chan Log.LogStuct, res chan bool, Ut Util.UtilStruct) (GRPCService, error) {
	grc := GRPCService{
		AddRequest: req,
		Response:   res,
		Ut: Ut,
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

	GRPCServicePackage.RegisterRPCServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(listener); err != nil {
  		log.Fatalf("failed to serve: %v", err)
 	}

}

func (Grc *GRPCService) RecurAddLog() {

}
