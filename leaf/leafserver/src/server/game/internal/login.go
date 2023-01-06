package internal

import (
	"leafserver/src/server/msg"
	"storage"

	"github.com/name5566/leaf/log"
)

func loginRequst(args []interface{}) {
	req := args[0].(*msg.LoginRequest)
	loginHandle(req)
}

func loginHandle(member *msg.LoginRequest) {
	if MD == nil {
		return
	}
	user, ok := Users[member.Uuid]
	if ok {
		user.Data = &storage.UserStorage{
			Uid:      member.Uuid,
			Nickname: member.Nickname,
			Icon:     member.Icon,
			Online:   1,
			Forever:  make(map[string]string),
		}
		r, e := user.Query()
		if e != nil {
			log.Debug("找不到用户，新建%v", e)
			user.Save()
		} else {
			// 更新新的名字和头像
			r.Nickname = member.Nickname
			r.Icon = member.Icon
			user.Data = r
		}
	} else {
		log.Debug("找不到用户")
	}
}

func offlineHanlde(uuid string) {
	if MD == nil {
		return
	}
	user, ok := Users[uuid]
	if ok {
		user.Update()
	}
}
