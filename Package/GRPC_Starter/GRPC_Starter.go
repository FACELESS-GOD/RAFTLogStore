package GRPC_Starter

import (
	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type GRPCServiceInterface interface {
	AddLog(Log.LogStuct) (bool, error)
}

type GRPCService struct {
	AddRequest chan Log.LogStuct
	Response   chan bool
}

func NewGRPCService(req chan Log.LogStuct, res chan bool) (GRPCService, error) {
	grc := GRPCService{
		AddRequest: req,
		Response:   res,
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

}

func (Grc *GRPCService) RecurAddLog() {

}
