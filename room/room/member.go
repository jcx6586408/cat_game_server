package room

import (
	"catLog"
	"proto/msg"
)

func NewMember(uuid string, uid string, nickname string, icon string, isMaster bool, isRobot bool) *msg.Member {
	if len(Members) > 0 {
		rm := Members[len(Members)-1]
		rm.Uuid = uuid
		rm.Uid = uid
		rm.Nickname = nickname
		rm.Icon = icon
		rm.IsMaster = isMaster
		rm.IsRobot = isRobot
		catLog.Log("重用成员")
		Delete(Members, rm)
		return rm
	}
	m := &msg.Member{}
	m.Uuid = uuid
	m.Uid = uid
	m.Nickname = nickname
	m.Icon = icon
	m.IsMaster = isMaster
	m.IsRobot = isRobot
	catLog.Log("新成员")
	return m
}

var Members = []*msg.Member{}

// 回收成员
func RecycleMember(member *msg.Member) {
	Members = append(Members, member)
	catLog.Log("成员回收", len(Members))
}

func Delete(a []*msg.Member, elem *msg.Member) []*msg.Member {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
		}
	}
	return a
}
