package Router

import (
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(Ctrl Controller.ControllerStruct) *gin.Engine {

	router := gin.Default()

	router.POST("/Add", Ctrl.AddLog)

	router.GET("/Get", Ctrl.GetLog)

	return router
}
