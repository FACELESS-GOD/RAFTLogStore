package GRPC_Starter

import (
	"context"
	"sync"
	"time"

	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type Server struct {
	mu sync.Mutex
	Ut Util.UtilStruct
	GRPCServicePackage.UnimplementedRPCServiceServer
}

func (Ser *Server) AppendRPC(Ctx context.Context, Req *GRPCServicePackage.AddLogRequest) (*GRPCServicePackage.AddLogResponse, error) {
	Ser.Ut.LastTouch = time.Now()
	if Req.IsHeartbeat == true {
		res := GRPCServicePackage.AddLogResponse{
			IsAnyError: false,
		}
		return &res, nil
	} else {

		for i := 0; i < len(Req.Log); i++ {
			tobeInjectedLog := Req.Log[i]

		}

		res := GRPCServicePackage.AddLogResponse{
			IsAnyError: false,
		}
		return &res, nil
	}
}

func (Ser *Server) RequestVoteRPC(Ctx context.Context, Req *GRPCServicePackage.RequestLogRequest) (*GRPCServicePackage.RequestLogResponse, error) {
	if Ser.Ut.Is_Voted != true {
		if Req.TermId > Ser.Ut.Term {
			Ser.Ut.Is_Voted = true
			res := GRPCServicePackage.RequestLogResponse{
				Vote: true,
			}
			return &res, nil
		} else if Req.TermId == Ser.Ut.Term {

			if Req.LogId > Ser.Ut.LogId {
				Ser.Ut.Is_Voted = true
				res := GRPCServicePackage.RequestLogResponse{
					Vote: true,
				}
				return &res, nil
			} else if Req.LogId == Ser.Ut.LogId {
				res := GRPCServicePackage.RequestLogResponse{
					Vote: false,
				}
				return &res, nil
			} else {
				res := GRPCServicePackage.RequestLogResponse{
					Vote: false,
				}
				return &res, nil
			}

		} else {
			res := GRPCServicePackage.RequestLogResponse{
				Vote: false,
			}
			return &res, nil
		}
	} else {
		res := GRPCServicePackage.RequestLogResponse{
			Vote: false,
		}
		return &res, nil
	}

}
