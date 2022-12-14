package internal

import (
	"leafserver/src/server/msg"
	pmsg "proto/msg"
	"remotemsg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"google.golang.org/grpc"
)

var Conn *grpc.ClientConn
var InnerClient pmsg.RoomClient

func RoomInit() {

}

func send(subRoom *Room, a gate.Agent, m *pmsg.Member, msgID int) {
	a.WriteMsg(&pmsg.RoomInfoReply{
		RoomID:         int32(subRoom.ID),
		PrepareMembers: subRoom.PrepareMembers,
		PlayingMembers: subRoom.PlayingMembers,
		Progress:       int32(subRoom.Cur),
		TotolQuestion:  int32(subRoom.QuestionCount),
		CurQuestion:    int32(subRoom.GetProgress()) + 1,
		ChangeMemeber:  m,
		MsgID:          int32(msgID),
		Question:       subRoom.GetQuestion(),
		ToTalTime:      int32(subRoom.GetPlayTime()),
		MaxMemeber:     int32(subRoom.MaxMember),
	})
}

// 房间创建
func roomCreate(args []interface{}) {
	req := args[0].(*pmsg.CreateRoomRequest)
	room := Manager.CreateRoom(req.Member)
	log.Debug("房间创建消息-------------------房间ID: %v---------------%v", room.ID, req.Member.Uuid)
	a := args[1].(gate.Agent)
	a.WriteMsg(&pmsg.CreateRoomReply{RoomID: int32(room.ID)})
}

func roomInfoGet(args []interface{}) {
	req := args[0].(*pmsg.CreateRoomRequest)
	room := Manager.CreateRoom(req.Member)
	log.Debug("房间获取消息-------------------房间ID: %v---------------%v", room.ID, req.Member.Uuid)
	a := args[1].(gate.Agent)
	send(room, a, nil, remotemsg.ROOMGET)
}

// 房间开始游戏
func roomStartPlay(args []interface{}) {
	req := args[0].(*msg.RoomStartPlay)
	log.Debug("房间开始游戏消息-------------------房间ID: %v", req.RoomID)
	Manager.StartPlay(req.RoomID)
}

// 加入准备房间
func roomAdd(args []interface{}) {
	req := args[0].(*pmsg.AddRequest)
	a := args[1].(gate.Agent)
	log.Debug("房间加入消息----------------房间ID: %v------------------%v", req.RoomID, req.Member.Uuid)
	_, code, err := Manager.AddFriendMember(int(req.RoomID), req.Member)
	if err == nil {
		// room.Send(remotemsg.ROOMADD)
	} else {
		log.Debug("返回加入房间错误码：%v", code)
		a.WriteMsg(&pmsg.RoomAddFail{
			Code: int32(code),
		})
	}
}

// 离开准备房间
func roomLeave(args []interface{}) {
	req := args[0].(*pmsg.LeaveRequest)
	log.Debug("房间离开消息---------------------房间ID: %v-------------%v", req.RoomID, req.Member.Uuid)
	// room, err := Manager.LeavePrepareMemeber(int(req.RoomID), req.Member)
	// if err == nil {
	// 	log.Debug("广播玩家离开的消息, %v")
	// 	room.Send(remotemsg.ROOMLEAVE)
	// }
}

// 答题
func roomAnswer(args []interface{}) {
	req := args[0].(*pmsg.Answer)
	log.Debug("收到答题消息------------------------%v", req)
	// Manager.AnswerQuestion(req)
}

// 房间解散
func roomOver(args []interface{}) {
	req := args[0].(*pmsg.OverRequest)
	log.Debug("房间解散消息----------------------------------", req.Member.Uuid)
	// _, err := Manager.OverRoom(int(req.RoomID))
	// if err == nil {
	// 	a := args[1].(gate.Agent)
	// 	a.WriteMsg(&pmsg.OverReply{State: true})
	// }
}

// 房间匹配
func roomMatch(args []interface{}) {
	req := args[0].(*pmsg.MatchRoomRequest)
	log.Debug("房间匹配消息----------------------------------%v", req.RoomID)
	roomID := int(req.RoomID)
	// room, err := Manager.GetRoomByID(roomID)
	// if err == nil {
	// 	room.Send(remotemsg.ROOMMATCHROOM)
	// 	Manager.MatchRoom(room)
	// }
}

// 个人匹配
func roomMatchMember(args []interface{}) {
	req := args[0].(*pmsg.MatchMemberRequest)
	log.Debug("房间个人消息----------------------------------%v", req.Member.Uuid)
	// room, err := Manager.MatchMember(req.Member)
	// if err == nil {
	// 	room.Send(remotemsg.ROOMMATCH)
	// }
}

// 房间匹配取消
func roomMatchCanel(args []interface{}) {
	req := args[0].(*pmsg.MatchRoomRequest)
	// _, err := Manager.MatchRoomCancel(int(req.RoomID))
	// if err == nil {
	// 	a := args[1].(gate.Agent)
	// 	a.WriteMsg(&pmsg.MatchRoomCancelReply{State: true})
	// }
}

// 个人匹配取消
func roomMatchMemberCanel(args []interface{}) {
	req := args[0].(*pmsg.LeaveRequest)
	// _, err := Manager.MatchMemberCancel(req)
	// if err == nil {
	// 	a := args[1].(gate.Agent)
	// 	a.WriteMsg(&pmsg.MatchMemberCancelReply{State: true})
	// }
}

// 成员复活
func roomMatchMemberRelive(args []interface{}) {
	req := args[0].(*pmsg.MemberReliveRequest)
	log.Debug("复活消息----------------: %v", req.Uuid)
	// room, err := Manager.Relive(req)
	// if err == nil {
	// 	room.SendRelive(req)
	// }
}
