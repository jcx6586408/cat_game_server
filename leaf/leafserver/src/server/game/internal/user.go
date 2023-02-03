package internal

import (
	"errors"
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

	DBNAME string

	COLLECT string
)

func MongoConnect() {
	c, err := mongodb.Dial(ServerConf.MongoDB.Url, ServerConf.MongoDB.SessionNum)
	if err != nil {
		return
	}
	// c.Close()
	DBNAME = ServerConf.MongoDB.DB
	COLLECT = ServerConf.MongoDB.Collection
	log.Release("芒果数据库连接成功********************")
	MD = c
}

func NewUserStorageData(uid, nickname, icon string) *storage.UserStorage {
	r := &storage.UserStorage{}
	r.Uid = uid
	r.Nickname = nickname
	r.Icon = icon
	r.Forever = make(map[string]string)
	return r
}

func (u *User) Save() {
	if MD == nil {
		return
	}
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	e := c.Insert(u.Data)
	if e != nil {
		log.Debug("数据插入失败%v", e)
	}
	MD.UnRef(s)
}

func (u *User) Update() error {
	if MD == nil {
		return errors.New("")
	}
	if u.Data == nil {
		return errors.New("")
	}
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	selector := bson.M{"uid": u.Data.Uid}
	log.Debug("离线保存数据=============: %v", u.Data.Uid)
	update := bson.M{"$set": bson.M{"uid": u.Data.Uid, "icon": u.Data.Icon, "nickname": u.Data.Nickname, "online": 0, "forever": u.Data.Forever}}
	if err := c.Update(selector, update); err != nil {
		MD.UnRef(s)
		return err
	}
	MD.UnRef(s)
	return nil
}

func (u *User) UpdateMap(key, value string) {
	if MD == nil {
		return
	}
	if u.Data == nil {
		return
	}
	if u.Data.Forever == nil {
		u.Data.Forever = make(map[string]string)
	}
	log.Debug("存储数据***********: %v|%v", key, value)
	u.Data.Forever[key] = value
}

func (u *User) GetData(key string) string {
	// 数据库不存在，直接返回空
	if MD == nil {
		return ""
	}
	if u.Data == nil {
		return ""
	}
	if u.Data.Forever == nil {
		u.Data.Forever = make(map[string]string)
	}
	val, ok := u.Data.Forever[key]
	log.Debug("获取数据----------: %v|%v", key, val)
	if !ok {
		return ""
	}
	return val
}

func (u *User) Query() (*storage.UserStorage, error) {
	if MD == nil {
		return nil, errors.New("")
	}
	var s = MD.Ref()
	var c = s.DB(DBNAME).C(COLLECT)
	var result *storage.UserStorage
	if err := c.Find(bson.M{"uid": u.Data.Uid}).One(&result); err != nil {
		MD.UnRef(s)
		log.Debug("查找玩家数据失败----: %v", u.Data.Uid)
		return nil, err
	}
	MD.UnRef(s)
	log.Debug("查找玩家数据++++: %v", u.Data)
	return result, nil
}
