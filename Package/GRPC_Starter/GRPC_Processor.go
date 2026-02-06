package GRPC_Starter

import (
	"context"
	"sync"
	"time"

	GRPCServicePackage "github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Package/GRPC_Mapper"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type Server struct {
	mu  sync.Mutex
	Ut  Util.UtilStruct
	Mdl Model.ModelStuct
	GRPCServicePackage.UnimplementedRPCServiceServer
}

func (Ser *Server) AppendRPC(Ctx context.Context, Req *GRPCServicePackage.AddLogRequest) (*GRPCServicePackage.AddLogResponse, error) {
	Ser.mu.Lock()
	defer Ser.mu.Unlock()

	Ser.Ut.LastTouch = time.Now()
	if Req.IsHeartbeat == true {

		Ser.Ut.VotedForWhom = Req.TermId

		Ser.Ut.Term = Req.TermId

		Ser.Ut.Is_Voted = false

		res := GRPCServicePackage.AddLogResponse{
			IsAnyError: false,
		}
		return &res, nil
	} else {

		for i := 0; i < len(Req.Log); i++ {
			tobeInjectedLog := Req.Log[i]
			if tobeInjectedLog.LogId < int32(len(Ser.Mdl.Arr)) {
				Ser.Mdl.Arr[int32(tobeInjectedLog.LogId)] = tobeInjectedLog.Text
			} else {
				Ser.Mdl.Arr = append(Ser.Mdl.Arr, tobeInjectedLog.Text)
				Ser.Mdl.Index = len(Ser.Mdl.Arr) - 1
				Ser.Ut.LogId = int32(len(Ser.Mdl.Arr) - 1)
			}

		}
		Ser.mu.Unlock()

		res := GRPCServicePackage.AddLogResponse{
			IsAnyError: false,
		}
		return &res, nil
	}
}

func (Ser *Server) RequestVoteRPC(Ctx context.Context, Req *GRPCServicePackage.RequestLogRequest) (*GRPCServicePackage.RequestLogResponse, error) {
	Ser.mu.Lock()
	defer Ser.mu.Unlock()

	if Ser.Ut.Is_Voted != true {
		if Req.TermId > Ser.Ut.Term {
			Ser.Ut.Term = Req.TermId

			Ser.Ut.Is_Voted = true

			Ser.Ut.VotedForWhom = Req.TermId

			res := GRPCServicePackage.RequestLogResponse{
				Vote:   true,
				LogId:  Ser.Ut.LogId,
				TermId: Ser.Ut.Term,
			}
			return &res, nil
		} else if Req.TermId == Ser.Ut.Term {

			if Req.LogId > Ser.Ut.LogId {

				Ser.Ut.Is_Voted = true

				Ser.Ut.VotedForWhom = Req.TermId

				res := GRPCServicePackage.RequestLogResponse{
					Vote:   true,
					LogId:  Ser.Ut.LogId,
					TermId: Ser.Ut.Term,
				}

				return &res, nil
			} else if Req.LogId == Ser.Ut.LogId {
				res := GRPCServicePackage.RequestLogResponse{
					Vote:   false,
					LogId:  Ser.Ut.LogId,
					TermId: Ser.Ut.Term,
				}
				return &res, nil
			} else {

				res := GRPCServicePackage.RequestLogResponse{
					Vote:   false,
					LogId:  Ser.Ut.LogId,
					TermId: Ser.Ut.Term,
				}
				return &res, nil
			}

		} else {
			res := GRPCServicePackage.RequestLogResponse{
				Vote:   false,
				LogId:  Ser.Ut.LogId,
				TermId: Ser.Ut.Term,
			}
			return &res, nil
		}
	} else {
		res := GRPCServicePackage.RequestLogResponse{
			Vote:   false,
			LogId:  Ser.Ut.LogId,
			TermId: Ser.Ut.Term,
		}
		return &res, nil
	}

}
