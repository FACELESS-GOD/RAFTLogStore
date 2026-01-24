package Model

import (
	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/GRPC_Starter"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type ModelInterFace interface {
	AddLog(Log.LogStuct) (bool, error)
	GetLog(int) (Log.LogStuct, error)
}

type ModelStuct struct {
	Utility     Util.UtilStruct
	Arr         []string
	GrpcService GRPC_Starter.GRPCServiceInterface
}

func NewModel(UT Util.UtilStruct, GrpcService GRPC_Starter.GRPCServiceInterface) (ModelStuct, error) {
	arr := make([]string, 10, 10)
	mdl := ModelStuct{Utility: UT, Arr: arr, GrpcService: GrpcService}
	return mdl, nil
}

func (Mdl *ModelStuct) AddLog(LogStuct Log.LogStuct) (bool, error) {

	if IsAdded, err := Mdl.GrpcService.AddLog(LogStuct) ; err != nil  {
		return false , err 		
	} else if IsAdded != true {
		return false , nil 		
	} else {
		Mdl.Arr = append(Mdl.Arr, LogStuct.Text)
		return true, nil
	}	
}

func (Mdl *ModelStuct) GetLog(Id int) (Log.LogStuct, error) {
	log := Log.LogStuct{}
	if Id <= len(Mdl.Arr) {
		log.Text = Mdl.Arr[Id]
	}
	return log, nil
}
