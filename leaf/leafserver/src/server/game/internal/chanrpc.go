package internal

import (
	"github.com/google/uuid"
	"github.com/name5566/leaf/gate"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("Login", login)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	guid := uuid.New().String()
	u := &User{
		Uuid:  guid,
		Agent: a,
	}
	AddUser(guid, u)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	DeleUser(a)
}
