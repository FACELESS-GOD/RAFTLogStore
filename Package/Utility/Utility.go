package Util

import (
	"errors"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/ServerMode"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
)

type UtilStruct struct {
	Mode        int
	ServerMode  int
	Routing_URL string
	Term        int
}

func NewUtil(Mode, Server_Mode int) (UtilStruct, error) {
	util := UtilStruct{}

	if Mode != State.Leader && Mode != State.Candidate && Mode != State.Follower {
		return util, errors.New("Invalid Mode")
	}

	if Server_Mode != ServerMode.Prod && Server_Mode != ServerMode.QA && Server_Mode != ServerMode.Test {
		return util, errors.New("Invalid Server_Mode")
	}

	util.Mode = Mode
	util.ServerMode = Server_Mode
	util.Term = 1
	return util, nil
}
