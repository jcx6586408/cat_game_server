package internal

import (
	"leafserver/src/server/msg"
	"storage"

	"github.com/name5566/leaf/log"
)

func loginRequst(args []interface{}) {
	req := args[0].(*msg.LoginRequest)
	log.Debug("玩家登录请求: uuid: %v", req.Uuid)
	loginHandle(req)
}

func loginHandle(member *msg.LoginRequest) {
	if MD == nil {
		return
	}
	user, ok := Users[member.Uuid]
	if ok {
		user.Data = &storage.UserStorage{
			Uid:      member.Uid,
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
			log.Debug("玩家登录时的数据: %v", user.Data)
		}
	} else {
		log.Debug("找不到用户")
		// user.Data = NewUserStorageData(member.Uid, member.Nickname, member.Icon)
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
