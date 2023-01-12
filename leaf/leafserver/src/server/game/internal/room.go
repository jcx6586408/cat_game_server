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
	ResetToPrepare()                     // 回归准备状态
	GetMembers() []*pmsg.Member          // 获取所有成员
	GetMemeber(uuid string) *pmsg.Member // 获取单个成员
	SendLeave(member *pmsg.Member)       // 发送玩家离开
	OnEndPlay()                          // 游戏结束处理
	Send(msgID int, change *pmsg.Member)
	Relive(uuid string)
	Answer(a *pmsg.Answer)     // 答题
	OfflinHanlder(uuid string) // 离线检测
	MemberReady(uuid string) bool
	MemberReadyCancel(uuid string) bool
	MemberLevelChange(uuid string, level int32)
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
	var members = []*pmsg.Member{}
	for _, v := range r.Members {
		if v.State == int32(MEMBERPLAYING) || v.State == int32(MEMBERMATCHING) {
			members = append(members, v)
		}
	}
	return members
}

func (r *Room) GetPrepareMembers() []*pmsg.Member {
	var members = []*pmsg.Member{}
	for _, v := range r.Members {
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			members = append(members, v)
		}
	}
	return members
}

func (r *Room) GetPlayingMembers() []*pmsg.Member {
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
	r.ChangeMemberState(MEMBERMATCHING)
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
	member.RoomID = int32(r.GetID())
	member.IsMaster = (r.GetMemberCount() <= 0) // 第一个人设置为房主
	member.State = int32(MEMEBERPREPARE)
	r.Members = append(r.Members, member)
	if r.BattleRoom != nil {
		log.Debug("加入新成员*************3")
		r.BattleRoom.Send(remotemsg.ROOMADD, member)
		a := Users[member.Uuid]
		if a != nil {
			a.WriteMsg(&pmsg.RoomInfoReply{
				RoomID:         int32(r.GetID()),
				PrepareMembers: r.GetPrepareMembers(),
				PlayingMembers: r.GetPlayingMembers(),
				Progress:       0,
				TotolQuestion:  0,
				CurQuestion:    0,
				ChangeMemeber:  member,
				MsgID:          int32(remotemsg.ROOMADD),
				Question:       nil,
				ToTalTime:      0,
				MaxMemeber:     int32(r.Max),
				BattleRoomID:   0,
			})
		}
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
		if v.State != int32(MEMEBENONERPREPARE) {
			v.State = int32(state)
			v.IsDead = false
		}
	}
}

func (r *Room) ResetToPrepare() {
	for _, v := range r.Members {
		if v.State != int32(MEMEBENONERPREPARE) {
			v.State = int32(MEMEBERPREPARE)
			v.IsDead = false
		}
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
	} else {
		log.Debug("战斗房间为, 复活失败, 战斗房间ID")
	}
}

func (m *Room) MemberReady(uuid string) bool {
	for _, v := range m.GetMembers() {
		if v.Uuid == uuid {
			if v.State == int32(MEMBERPLAYING) || v.State == int32(MEMBERMATCHING) {
				log.Debug("ready成员正则匹配或游玩，无法改变准备状态")
				return false
			}
			log.Debug("ready改变状态成功为准备状态")
			v.State = int32(MEMEBERPREPARE)
			return true
		}
	}
	return false
}

func (m *Room) MemberReadyCancel(uuid string) bool {
	for _, v := range m.GetMembers() {
		if v.Uuid == uuid {
			if v.State == int32(MEMBERPLAYING) || v.State == int32(MEMBERMATCHING) {
				log.Debug("readyCancel成员正则匹配或游玩，无法改变准备状态")
				return false
			}
			log.Debug("readyCancel改变状态成功为非准备状态")
			v.State = int32(MEMEBENONERPREPARE)
			return true
		}
	}
	return false
}

func (m *Room) MemberLevelChange(uuid string, level int32) {
	for _, v := range m.GetMembers() {
		if v.Uuid == uuid {
			log.Debug("修改等级成功: 原等级:%v  新等级:%v", v.Level, level)
			v.Level = level
			break
		}
	}
}
