package internal

import (
	"storage"

	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"github.com/siddontang/go/bson"
)

type User struct {
	Uuid string
	Data *storage.UserStorage // 数据库数据
	gate.Agent
}

var (
	Users = make(map[string]*User)

	AgentUsers = make(map[gate.Agent]string)

	DBNAME = "sheep"

	COLLECT = "user"
)

func MongoConnect() {
	c, err := mongodb.Dial("localhost:27017", 100)
	if err != nil {
		panic(err)
	}
	// c.Close()
	log.Debug("芒果数据库连接成功********************")
	MD = c
}

func (u *User) Save() {
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	e := c.Insert(u.Data)
	if e != nil {
		log.Debug("数据插入失败%v", e)
	}
	MD.UnRef(s)
}

func (u *User) Update() error {
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	selector := bson.M{"uid": u.Data.Uid}
	update := bson.M{"$set": bson.M{"uid": u.Data.Uid, "icon": u.Data.Icon, "nickname": u.Data.Nickname, "online": 1, "forever": u.Data.Forever}}
	if err := c.Update(selector, update); err != nil {
		MD.UnRef(s)
		return err
	}
	MD.UnRef(s)
	return nil
}

func (u *User) UpdateMap(key, value string) {
	if u.Data.Forever == nil {
		u.Data.Forever = make(map[string]string)
	}
	u.Data.Forever[key] = value
}

func (u *User) Query() (*storage.UserStorage, error) {
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	var result *storage.UserStorage
	if err := c.Find(bson.M{"uid": u.Data.Uid}).One(&result); err != nil {
		MD.UnRef(s)
		return nil, err
	}
	MD.UnRef(s)
	return result, nil
}
