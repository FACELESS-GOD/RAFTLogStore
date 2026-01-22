package Model

import (
	"testing"

	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
	suite.Suite
	Mdl ModelStuct
}

func TestMain(m *testing.T) {
	suite.Run(m, &TestStruct{})
}

func (Ts *TestStruct) SetupSuite() {
	util, err := Util.NewUtil(2, 1)

	if err != nil {
		Ts.FailNow(err.Error())
	}
	mdl, err := NewModel(util)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Mdl = mdl

}

func (Its *TestStruct) TestAddLog() {
	log := LogStuct{Text: "Hello World!"}
	IsAdded, err := Its.Mdl.AddLog(log)
	Its.Require().Nil(err)
	Its.Require().Equal(IsAdded == true, true)

}

func (Its *TestStruct) TestGetLog() {

	log := LogStuct{Text: "Hello World!"}
	IsAdded, err := Its.Mdl.AddLog(log)

	if err != nil {
		Its.FailNow(err.Error())
	}

	if IsAdded == false {
		Its.FailNow("Log not added Properly.")
	}

	newLogStruct, err := Its.Mdl.GetLog(len(Its.Mdl.Arr)-1)
	Its.Require().Nil(err)
	Its.Require().NotNil(newLogStruct.Text)
	Its.Require().Equal(len(newLogStruct.Text)>0 , true)

}
