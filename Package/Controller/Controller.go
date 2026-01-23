package Controller

import (
	"net/http"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
	"github.com/gin-gonic/gin"
)

type GetLogRequest struct {
	ID int `json:"Id"`
}

type AddLogResponse struct {
	IsAdded       bool
	IsAnyError    bool
	ErrorMessages []string
}

func NewAddLogResponse() AddLogResponse {
	arr := []string{}
	return AddLogResponse{IsAdded: false, IsAnyError: false, ErrorMessages: arr}
}

type ControllerStruct struct {
	Utility Util.UtilStruct
	Mdl     Model.ModelInterFace
}

func NewController(Util Util.UtilStruct, Mdl Model.ModelInterFace) (ControllerStruct, error) {
	ctrl := ControllerStruct{
		Utility: Util,
		Mdl:     Mdl,
	}
	return ctrl, nil
}

func (Ctx *ControllerStruct) AddLog(Gctx *gin.Context) {

	if Ctx.Utility.Mode == State.Follower || Ctx.Utility.Mode == State.Candidate {
		Gctx.Redirect(http.StatusFound, Ctx.Utility.Routing_URL)
		return
	}
	response := NewAddLogResponse()
	log := Model.LogStuct{}

	err := Gctx.ShouldBindBodyWithJSON(&log)

	if err != nil {
		response.IsAnyError = true
		response.IsAdded = false
		response.ErrorMessages = append(response.ErrorMessages, err.Error())
		Gctx.Status(http.StatusInternalServerError)
		return
	} else {

		isAdded, err := Ctx.Mdl.AddLog(log)

		if err != nil {
			response.IsAnyError = true
			response.IsAdded = false
			response.ErrorMessages = append(response.ErrorMessages, err.Error())
			Gctx.Status(http.StatusInternalServerError)
			return
		} else if isAdded != true {
			response.IsAnyError = true
			response.IsAdded = false
			response.ErrorMessages = append(response.ErrorMessages, "Error Has occured!")
			Gctx.Status(http.StatusInternalServerError)
			return
		} else {
			response.IsAnyError = false
			response.IsAdded = true
			Gctx.JSON(http.StatusOK, response)
			return
		}
	}

}

func (Ctx *ControllerStruct) GetLog(Gctx *gin.Context) {

	if Ctx.Utility.Mode == State.Follower || Ctx.Utility.Mode == State.Candidate {
		Gctx.Redirect(http.StatusFound, Ctx.Utility.Routing_URL)
		return
	}

	response := NewAddLogResponse()
	logIDStruct := GetLogRequest{}

	err := Gctx.ShouldBindBodyWithJSON(&logIDStruct)

	if err != nil {
		response.IsAnyError = true
		response.IsAdded = false
		response.ErrorMessages = append(response.ErrorMessages, err.Error())
		Gctx.Status(http.StatusInternalServerError)
		return

	} else {

		log, err := Ctx.Mdl.GetLog(logIDStruct.ID)

		if err != nil {
			response.IsAnyError = true
			response.IsAdded = false
			response.ErrorMessages = append(response.ErrorMessages, err.Error())
			Gctx.Status(http.StatusInternalServerError)
			return

		} else {
			response.IsAnyError = false
			response.IsAdded = true
			Gctx.JSON(http.StatusOK, log)
			return

		}
	}
}
