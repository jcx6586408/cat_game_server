package internal

import (
	"storage"

	"github.com/labstack/gommon/log"
)

func loginHandle(uuid, uid string) {
	log.Debug("登录处理*******************uuid: %s; uid(openid): %s", uuid, uid)
	user, ok := Users[uuid]
	if ok {
		u := storage.Queryuser(uid)
		if u != nil {
			user.Data = u
			storage.UpdateOnline(u, 1)
		} else {
			user.Data = &storage.UserStorage{
				Uid:    uid,
				Online: 1,
			}
			storage.Insectuser(user.Data)

		}
	}
}
