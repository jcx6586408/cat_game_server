package internal

import (
	"storage"

	"github.com/name5566/leaf/gate"
)

type User struct {
	Uuid string
	Data *storage.UserStorage // 数据库数据
	gate.Agent
}

var Users = make(map[string]*User)

var AgentUsers = make(map[gate.Agent]string)
