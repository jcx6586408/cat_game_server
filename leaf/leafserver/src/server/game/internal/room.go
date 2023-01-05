package internal

import (
	pmsg "proto/msg"
	"remotemsg"

	"github.com/name5566/leaf/log"
)

type Roomer interface {
	Matching()                           // 开始匹配
	MatchingCancel() bool                // 取消匹配
	AddMember(member *pmsg.Member) bool  // 加入成员
	LeaveMember(member *pmsg.Member)     // 离开成员
	ChangeMemberState(state int)         // 改变成员状态
	GetMembers() []*pmsg.Member          // 获取所有成员
	GetMemeber(uuid string) *pmsg.Member // 获取单个成员
	SendLeave(member *pmsg.Member)       // 发送玩家离开
	OnEndPlay()                          // 游戏结束处理
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
	r.BattleRoom = nil
}

func (r *Room) OnClose() {
	r.Members = nil
	r.BattleRoom = nil
}

func (r *Room) GetMembers() []*pmsg.Member {
	return r.Members
}

func (r *Room) GetMemeber(uuid string) *pmsg.Member {
	for _, v := range r.Members {
		if v.Uuid == uuid {
			return v
		}
	}
	return nil
}

func (r *Room) GetMaster() *pmsg.Member {
	for _, v := range r.Members {
		if v.IsMaster {
			return v
		}
	}
	return nil
}

func (r *Room) Matching() {
	master := r.GetMaster()
	// 游戏状态不运行匹配
	if master != nil && master.State == int32(MEMBERPLAYING) {
		return
	}
	br := BattleManager.MatchRoom(r)
	r.BattleRoom = br
	r.BattleRoom.Send(remotemsg.ROOMMATCHROOM, nil)
}

func (r *Room) MatchingCancel() bool {
	bo := BattleManager.MatchRoomzCancel(r)
	r.SendRoomMatchCancel(bo)
	return bo
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
	loginHandle(member) // 登录处理
	member.RoomID = int32(r.GetID())
	member.IsMaster = (r.GetMemberCount() <= 0) // 第一个人设置为房主
	member.State = int32(MEMEBERPREPARE)
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

func (r *Room) LeaveMember(member *pmsg.Member) {
	log.Debug("玩家离开房间, %v", member.GetUuid())
	member.RoomID = 0
	r.Members = r.delete(r.Members, member)
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
		r.SendLeave(member)
	}
	if len(r.Members) <= 0 {
		RoomManager.Destroy(r.ID)
	}
}

func (r *Room) ChangeMemberState(state int) {
	for _, v := range r.Members {
		v.State = int32(state)
		v.IsDead = false
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
		log.Debug("找到答题战斗房间")
		m.BattleRoom.Answer(a)
	}
}

func (m *Room) Relive(uuid string) {
	if m.BattleRoom != nil {
		m.BattleRoom.Relive(uuid)
	}
}
