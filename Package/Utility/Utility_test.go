package Util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
	suite.Suite
}

func TestMain(m *testing.T) {
	suite.Run(m, &TestStruct{})
}

func (Its *TestStruct) TestNewUtil() {

	util , err := NewUtil(0,0)
	Its.Require().NotNil(err)
	Its.Require().NotNil(util)

	util , err = NewUtil(0,1)
	Its.Require().NotNil(err)
	Its.Require().NotNil(util)

	util , err = NewUtil(1,0)
	Its.Require().NotNil(err)
	Its.Require().NotNil(util)

	util , err = NewUtil(2,1)
	Its.Require().Nil(err)
	Its.Require().NotNil(util)

}
