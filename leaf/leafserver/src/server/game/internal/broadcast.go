package room

import (
	pmsg "proto/msg"

	"github.com/name5566/leaf/log"
)

type Agent interface {
	WriteMsg(data interface{})
}

func (m *Room) sendbase(Users map[string]Agent, call func(Agent)) {
	for _, v := range m.Members {
		a := Users[v.Uuid]
		// 准备人员不通知
		if v.State == int32(MEMEBERPREPARE) {
			continue
		}
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

func (m *BattleRoom) sendbase(Users map[string]Agent, call func(Agent, Roomer, *pmsg.Member)) {
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
			a := Users[v.Uuid]
			// 准备人员不通知
			if v.State == int32(MEMEBERPREPARE) {
				continue
			}
			// 离线玩家不通知
			if a == nil {
				log.Debug("用户不存在%v", v.Uuid)
				continue
			}
			call(a, room, v)
		}
	}
}

// 离开
func (m *Room) SendLeave(Users map[string]Agent, msgID int, member *pmsg.Member) {
	m.sendbase(Users, func(a Agent) {

	})
}

// 战斗中离开
func (m *BattleRoom) SendLeave(Users map[string]Agent, msgID int, member *pmsg.Member) {
	m.sendbase(Users, func(a Agent, room Roomer, member *pmsg.Member) {
		room.SendLeave(Users, msgID, member)
	})
}

// 发送复活
func (m *BattleRoom) SendRelive(Users map[string]Agent, req *pmsg.MemberReliveRequest) {
	m.sendbase(Users, func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.MemberReliveReply{
			Uuid:   req.Uuid,
			Answer: member.Answer[m.LibAnswer.Progress],
		})
	})
}

// 发送答题
func (m *BattleRoom) SendAnswer(Users map[string]Agent, uuid string, qid int32, result string) {
	m.sendbase(Users, func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.Answer{
			Uuid:       uuid,
			RoomID:     int32(m.ID),
			QuestionID: qid,
			Result:     result,
		})
	})
}

// 发送计时
func (m *BattleRoom) SendTime(Users map[string]Agent, cur int) {
	m.sendbase(Users, func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.RoomTime{
			Time: int32(cur),
		})
	})
}

func (m *BattleRoom) send(Users map[string]Agent, msgID int, change *pmsg.Member) {
	m.sendbase(Users, func(a Agent, room Roomer, member *pmsg.Member) {
		a.WriteMsg(&pmsg.RoomInfoReply{
			RoomID:         int32(room.GetID()),
			PrepareMembers: nil,
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
