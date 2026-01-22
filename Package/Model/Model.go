package Model

import Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"

type ModelInterFace interface {
	AddLog(LogStuct) (bool, error)
	GetLog(int) (LogStuct, error)
}

type LogStuct struct {
	Text string `json:"Text"`
}

type ModelStuct struct {
	Utility Util.UtilStruct
	Arr     []string
}

func NewModel(UT Util.UtilStruct) (ModelStuct, error) {
	arr := make([]string, 10, 10)
	mdl := ModelStuct{Utility: UT, Arr: arr}
	return mdl, nil
}

func (Mdl *ModelStuct) AddLog(LogStuct LogStuct) (bool, error) {	
	Mdl.Arr = append(Mdl.Arr, LogStuct.Text)
	return true, nil
}

func (Mdl *ModelStuct) GetLog(Id int) (LogStuct, error) {
	log := LogStuct{}
	if Id < len(Mdl.Arr) {
		log.Text = Mdl.Arr[Id]
	}
	return log, nil
}
