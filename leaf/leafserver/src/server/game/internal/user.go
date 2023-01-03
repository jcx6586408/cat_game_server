package internal

import (
	"storage"

	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/gate"
)

type User struct {
	Uuid string
	Data *storage.UserStorage // 数据库数据
	gate.Agent
}

var Users = make(map[string]*User)

var AgentUsers = make(map[gate.Agent]string)

func MongoConnect() {
	c, err := mongodb.Dial(":27017", 10)
	if err != nil {
		panic(err)
	}
	// c.Close()
	MD = c
}

func (u *User) Save() {
	MD.Ref().DB("sheep").C("user").Insert()
}
