package internal

import (
	pmsg "proto/msg"
	"remotemsg"

	"github.com/name5566/leaf/log"
)

type Agent interface {
	WriteMsg(data interface{})
}

func (m *Room) sendbase(call func(Agent)) {
	for _, v := range m.Members {
		a := Users[v.Uuid]
		// 机器人不通知
		if v.IsRobot {
			continue
		}
		// 离线玩家不通知
		if a == nil {
			log.Debug("用户不存在%v", v.Uuid)
			continue
		}
		call(a)
	}
}

func (m *BattleRoom) sendbase(call func(Agent, Roomer, *pmsg.Member)) {
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
			a := Users[v.Uuid]
			// 离线玩家不通知
			if a == nil {
				log.Debug("用户不存在%v", v.Uuid)
				continue
			}
			call(a, room, v)
		}
	}
}

func (m *Room) SendRoomMatchCancel(bo bool) {
	m.sendbase(func(a Agent) {
		log.Debug("取消匹配发送************: %v", bo)
		a.WriteMsg(&pmsg.MatchRoomCancelReply{State: bo})
	})
}

// 离开
func (m *Room) SendLeave(member *pmsg.Member) {
	m.Send(remotemsg.ROOMLEAVE, member)
}

// 战斗中离开
func (m *BattleRoom) SendLeave(lmember *pmsg.Member) {
	m.sendbase(func(a Agent, room Roomer, member *pmsg.Member) {
		room.SendLeave(lmember)
	})
}

// 发送复活
func (m *BattleRoom) SendRelive(uuid string) {
	m.sendbase(func(a Agent, room Roomer, member *pmsg.Member) {
		log.Debug("%v复活答案: %v", uuid, member.Answer[m.LibAnswer.Progress].Result)
		target := room.GetMemeber(uuid)
		if target != nil {
			a.WriteMsg(&pmsg.MemberReliveReply{
				Uuid:   uuid,
				Answer: target.Answer[m.LibAnswer.Progress],
			})
		}
	})
}

// 发送答题
func (m *BattleRoom) SendAnswer(uuid string, qid int32, result string) {
	m.sendbase(func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.Answer{
			Uuid:       uuid,
			RoomID:     int32(m.ID),
			QuestionID: qid,
			Result:     result,
		})
	})
}

// 发送计时
func (m *BattleRoom) SendTime(cur int) {
	m.sendbase(func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.RoomTime{
			Time: int32(cur),
		})
	})
}

func (m *BattleRoom) Send(msgID int, change *pmsg.Member) {
	m.sendbase(func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.RoomInfoReply{
			RoomID:         int32(room.GetID()),
			PrepareMembers: m.GetPrepareMembers(),
			PlayingMembers: m.GetPlayingMembers(),
			Progress:       int32(m.Cur),
			TotolQuestion:  int32(m.QuestionCount),
			CurQuestion:    int32(m.GetProgress() + 1),
			ChangeMemeber:  change,
			MsgID:          int32(msgID),
			Question:       m.GetQuestion(),
			ToTalTime:      int32(m.GetPlayTime()),
			MaxMemeber:     int32(m.Max),
			BattleRoomID:   int32(m.ID),
		})
	})
}

func (m *Room) Send(msgID int, change *pmsg.Member) {
	m.sendbase(func(a Agent) {
		a.WriteMsg(&pmsg.RoomInfoReply{
			RoomID:         int32(m.GetID()),
			PrepareMembers: m.Members,
			PlayingMembers: nil,
			Progress:       0,
			TotolQuestion:  0,
			CurQuestion:    0,
			ChangeMemeber:  change,
			MsgID:          int32(msgID),
			Question:       nil,
			ToTalTime:      0,
			MaxMemeber:     int32(m.Max),
			BattleRoomID:   0,
		})
	})
}
