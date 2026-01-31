package ServiceRegistry

import Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"

type Service struct {
	Ut         Util.UtilStruct
	ServerList []string
}

func RegisterService(Util Util.UtilStruct) (Service, error) {
	ser := Service{Ut: Util, ServerList: make([]string, 10)}
	// Register the service and populate the pier []string in Utilstruct
	return ser, nil
}
