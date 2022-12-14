package internal

import (
	"leafserver/src/server/msg"

	"github.com/google/uuid"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
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
	Users[guid] = u
	// 反注册
	AgentUsers[a] = guid
	log.Debug("玩家登录--------------------uuid: %v", guid)
	// 下发uuid
	a.WriteMsg(&msg.Login{
		Uuid: guid,
	})
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	guid := AgentUsers[a]
	log.Debug("玩家离线--------------------uuid: %v", guid)
	// storage.OfflineHandle(Users[guid].Data) // 离线保存
	delete(Users, guid)
	delete(AgentUsers, a)
	Manager.OfflineMemeber(guid)
}
