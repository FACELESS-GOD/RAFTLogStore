package Model

import (
	"sync"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

type ModelInterFace interface {
	AddLog(Log.LogStuct) (bool, error)
	GetLog(int) (Log.LogStuct, error)
}

type ModelStuct struct {
	Utility    Util.UtilStruct
	Arr        []string
	AddLogChan chan Log.LogStuct
	mu         *sync.Mutex
	Index      int
}

func NewModel(UT Util.UtilStruct, Mu *sync.Mutex) (ModelStuct, error) {
	arr := make([]string, 10)
	mdl := ModelStuct{Utility: UT, Arr: arr, mu: Mu}
	return mdl, nil
}

func (Mdl *ModelStuct) AddLog(LogStuct Log.LogStuct) (bool, error) {
	Mdl.AddLogChan <- LogStuct
	Mdl.mu.Lock()
	Mdl.Arr = append(Mdl.Arr, LogStuct.Text)
	Mdl.mu.Unlock()
	return true, nil
}

func (Mdl *ModelStuct) GetLog(Id int) (Log.LogStuct, error) {
	log := Log.LogStuct{}
	if Id <= len(Mdl.Arr) {
		log.Text = Mdl.Arr[Id]
	}
	return log, nil
}