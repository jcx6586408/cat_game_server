package room

import (
	pmsg "proto/msg"
)

type Roomer interface {
	Matching()                                                        // 开始匹配
	MatchingCancel()                                                  // 取消匹配
	AddMember(member Memberer) bool                                   // 加入成员
	LeaveMember(member Memberer)                                      // 加入成员
	ChangeMemberState(state int)                                      // 改变成员状态
	GetMembers() []*pmsg.Member                                       // 获取所有成员
	SendLeave(Users map[string]Agent, msgID int, member *pmsg.Member) // 发送玩家离开
	Roombaseer
}

type Room struct {
	ID      int
	Members []*pmsg.Member
	Max     int
}

func (r *Room) GetID() int {
	return r.ID
}

func (r *Room) OnInit() {

}

func (r *Room) OnClose() {

}

func (r *Room) GetMembers() []*pmsg.Member {
	return r.Members
}

func (r *Room) Matching() {
	BattleManager.MatchRoom(r)
}

func (r *Room) MatchingCancel() {
	BattleManager.MatchRoomzCancel(r)
}

func (r *Room) AddMember(member Memberer) bool {
	if r.Max <= r.GetMemberCount() {
		return false
	}
	r.Members = append(r.Members, member.(*pmsg.Member))
	return true
}

func (r *Room) LeaveMember(member Memberer) {
	r.Members = r.delete(r.Members, member.(*pmsg.Member))
}

func (r *Room) ChangeMemberState(state int) {
	for _, v := range r.Members {
		v.State = int32(state)
	}
}

func (r *Room) Full() bool {
	return r.GetMemberCount() >= r.Max
}

func (r *Room) GetMemberCount() int {
	return len(r.Members)
}

func (m *Room) delete(a []*pmsg.Member, elem *pmsg.Member) []*pmsg.Member {
	for i := 0; i < len(a); i++ {
		if a[i].GetUuid() == elem.GetUuid() {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
