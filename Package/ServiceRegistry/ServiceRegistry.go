package ServiceRegistry

import (
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type Service struct {
	Ut         Util.UtilStruct
	ServerList []string
}

func RegisterService(Util Util.UtilStruct, Mdl Model.ModelStuct) (Service, error) {
	ser := Service{Ut: Util, ServerList: make([]string, 10)}
	// Register the service and populate the pier []string in Utilstruct

	ser.CatchUp(Mdl)
	return ser, nil
}

func (Service *Service) CatchUp(Mdl Model.ModelStuct) {
	// To Take Data From Snapshot.
}
