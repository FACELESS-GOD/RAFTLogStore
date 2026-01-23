package main

import (
	"log"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/ServerMode"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Controller"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Model"
	"github.com/FACELESS-GOD/RAFTLogStore/Package/Router"
	Util "github.com/FACELESS-GOD/RAFTLogStore/Package/Utility"
)

func main() {
	Util, err := Util.NewUtil(State.Leader, ServerMode.Test)

	if err != nil {
		log.Fatal(err.Error())
	}

	mdl, err := Model.NewModel(Util)

	if err != nil {
		log.Fatal(err.Error())
	}

	ctrl, err := Controller.NewController(Util, &mdl)

	router := Router.NewRouter(ctrl)

	router.Run("localhost:8080")

}
