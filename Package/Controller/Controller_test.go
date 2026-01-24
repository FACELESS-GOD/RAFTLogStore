package Controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	Log "github.com/FACELESS-GOD/RAFTLogStore/Helper/LogDescription"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TestGRPCStruct struct {
}

func (Ts *TestGRPCStruct) AddLog(Log.LogStuct) (bool, error) {
	return true, nil
}

type TestControllerStruct struct {
	suite.Suite
	Ctrl ControllerStruct
	Mdl  Model.ModelStuct
}

func (Its *TestControllerStruct) Reset() {

}

func TestMain(m *testing.T) {
	suite.Run(m, &TestControllerStruct{})
}

func (Ts *TestControllerStruct) SetupSuite() {
	util, err := Util.NewUtil(1, 1)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	grpcStruct := TestGRPCStruct{}
	mdl, err := Model.NewModel(util, &grpcStruct)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Mdl = mdl

	ctrl, err := NewController(util, &mdl)

	if err != nil {
		Ts.FailNow(err.Error())
	}

	Ts.Ctrl = ctrl

}

func (Its *TestControllerStruct) TestAddLog() {
	log := Log.LogStuct{Text: "Hello World!"}

	jsonData, err := json.Marshal(&log)

	if err != nil {
		Its.FailNow(err.Error())
	}

	router := gin.Default()

	router.POST("/Add", Its.Ctrl.AddLog)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/Add", strings.NewReader(string(jsonData)))

	if err != nil {
		Its.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req)

	Its.Require().Equal(200, recorder.Code)

}

func (Its *TestControllerStruct) TestGetLog() {
	LogId := GetLogRequest{ID: len(Its.Mdl.Arr) - 1}

	jsonData, err := json.Marshal(&LogId)

	if err != nil {
		Its.FailNow(err.Error())
	}

	router := gin.Default()

	router.GET("/Get", Its.Ctrl.GetLog)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/Get", strings.NewReader(string(jsonData)))

	if err != nil {
		Its.FailNow(err.Error())
	}

	router.ServeHTTP(recorder, req)

	Its.Require().Equal(200, recorder.Code)

	obj := json.NewDecoder(recorder.Body)

	fmt.Println(obj)
}

func (Its *TestControllerStruct) BeforeTest(SuiteName string, TestName string) {
	switch TestName {
	case "TestGetLog":
		Its.InjectLog()
	}
}

func (Its *TestControllerStruct) InjectLog() {
	log := Log.LogStuct{"New Hello World!"}
	Its.Mdl.Arr = append(Its.Mdl.Arr, log.Text)
}
