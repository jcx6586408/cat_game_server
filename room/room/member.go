package room

import (
	"proto/msg"
)

func NewMember(uuid string, uid string, nickname string, icon string, isMaster bool, isRobot bool) *msg.Member {
	m := &msg.Member{}
	m.Uuid = uuid
	m.Uid = uid
	m.Nickname = nickname
	m.Icon = icon
	m.IsMaster = isMaster
	m.IsRobot = isRobot
	return m
}
