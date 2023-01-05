package internal

import (
	"storage"

	"github.com/name5566/leaf/log"
)

func loginHandle(member Memberer) {
	user, ok := Users[member.GetUuid()]
	if ok {
		user.Data = &storage.UserStorage{
			Uid:      member.GetUid(),
			Nickname: member.GetNickname(),
			Icon:     member.GetIcon(),
			Online:   1,
			Forever:  make(map[string]string),
		}
		r, e := user.Query()
		if e != nil {
			log.Debug("找不到用户，新建%v", e)
			user.Save()
		} else {
			user.Data = r
		}
	} else {
		log.Debug("找不到用户")
	}
}

func offlineHanlde(uuid string) {
	user, ok := Users[uuid]
	if ok {
		user.Update()
	}
}
