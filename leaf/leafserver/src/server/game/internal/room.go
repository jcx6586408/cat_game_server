package internal

import (
	pmsg "proto/msg"
	"remotemsg"

	"github.com/name5566/leaf/log"
)

type Roomer interface {
	Matching()                          // 开始匹配
	MatchingCancel()                    // 取消匹配
	AddMember(member *pmsg.Member) bool // 加入成员
	LeaveMember(member Memberer)        // 离开成员
	ChangeMemberState(state int)        // 改变成员状态
	GetMembers() []*pmsg.Member         // 获取所有成员
	SendLeave(member *pmsg.Member)      // 发送玩家离开
	OnEndPlay()                         // 游戏结束处理
	Send(msgID int, change *pmsg.Member)
	Relive(uuid string)
	Answer(a *pmsg.Answer)     // 答题
	OfflinHanlder(uuid string) // 离线检测
	Roombaseer
}

type Room struct {
	ID         int
	Members    []*pmsg.Member
	Max        int
	BattleRoom BattleRoomer
}

func (r *Room) GetID() int {
	return r.ID
}

func (r *Room) OnInit() {
	r.Members = []*pmsg.Member{}
	r.Max = RoomConf.MaxInvite
}

func (r *Room) OnClose() {

}

func (r *Room) GetMembers() []*pmsg.Member {
	return r.Members
}

func (r *Room) Matching() {
	br := BattleManager.MatchRoom(r)
	r.BattleRoom = br
	log.Debug("广播315房间匹配成功的消息--------------------------")
	r.BattleRoom.Send(remotemsg.ROOMMATCHROOM, nil)

}

func (r *Room) MatchingCancel() {
	BattleManager.MatchRoomzCancel(r)
	r.BattleRoom = nil
}

func (r *Room) OfflinHanlder(uuid string) {
	for _, v := range r.Members {
		if v.Uuid == uuid {
			r.LeaveMember(v)
		}
	}
}

func (r *Room) OnEndPlay() {
	r.ChangeMemberState(MEMEBERPREPARE)
	r.BattleRoom = nil
	log.Debug("清除战斗房间========================")
}

func (r *Room) AddMember(member *pmsg.Member) bool {
	log.Debug("加入新成员*************1")
	if r.Max <= r.GetMemberCount() {
		return false
	}
	log.Debug("加入新成员*************2")
	member.IsMaster = (r.GetMemberCount() <= 0) // 第一个人设置为房主
	r.Members = append(r.Members, member)
	if r.BattleRoom != nil {
		log.Debug("加入新成员*************3")
		r.BattleRoom.Send(remotemsg.ROOMADD, member)
	} else {
		log.Debug("加入新成员*************4")
		r.Send(remotemsg.ROOMADD, member) // 广播加入房间
	}
	return true
}

func (r *Room) LeaveMember(member Memberer) {
	log.Debug("玩家离开房间, %v", member.GetUuid())
	r.Members = r.delete(r.Members, member.(*pmsg.Member))
	// 如果离开的人是房主，则进行房主转移
	if member.GetIsMaster() {
		if r.GetMemberCount() > 0 {
			log.Debug("转移房主, %v", member.GetUuid())
			r.Members[0].IsMaster = true
		}
	}
	if r.BattleRoom != nil {
		// 检查是否有人
		log.Debug("玩家离开房间===============1, %v", member.GetUuid())
		r.BattleRoom.OnLeave(r, member)
	} else {
		log.Debug("玩家离开房间*****************2, %v", member.GetUuid())
		r.SendLeave(member.(*pmsg.Member))
	}
	if len(r.Members) <= 0 {
		RoomManager.Destroy(r.ID)
	}
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

func (m *Room) Answer(a *pmsg.Answer) {
	if m.BattleRoom != nil {
		m.BattleRoom.Answer(a)
	}
}

func (m *Room) Relive(uuid string) {
	if m.BattleRoom != nil {
		m.BattleRoom.Relive(uuid)
	}
}
