package Util

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/FACELESS-GOD/RAFTLogStore/Helper/ServerMode"
	"github.com/FACELESS-GOD/RAFTLogStore/Helper/State"
)

type UtilStruct struct {
	Mode            int
	ServerMode      int
	Routing_URL     string
	Term            int32
	ElectionTimeout time.Duration
	pier            []string
	LastTouch       time.Time
	Is_Voted        bool
	LogId           int32
	Mu              *sync.Mutex
}

func NewUtil(Mode, Server_Mode int, Mu *sync.Mutex) (UtilStruct, error) {
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
	offcet := 100 * rand.Int()
	util.ElectionTimeout = time.Millisecond * time.Duration(offcet)
	util.Is_Voted = true
	util.Mu = Mu
	return util, nil
}
