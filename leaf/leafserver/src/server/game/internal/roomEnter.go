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

// 房间创建
func roomCreate(args []interface{}) {
	req := args[0].(*pmsg.CreateRoomRequest)
	room := manager.Create()
	room.AddMember(req.Member)
	a := args[1].(gate.Agent)
	log.Debug("房间创建消息-------------------房间ID: %v---------------%v", room.GetID(), req.Member.Uuid)
	a.WriteMsg(&pmsg.CreateRoomReply{RoomID: int32(room.GetID())})
}

func roomInfoGet(args []interface{}) {
	// req := args[0].(*pmsg.CreateRoomRequest)
	// room := Manager.CreateRoom(req.Member)
	// log.Debug("房间获取消息-------------------房间ID: %v---------------%v", room.ID, req.Member.Uuid)
	// a := args[1].(gate.Agent)
	// send(room, a, nil, remotemsg.ROOMGET)
}

// 房间开始游戏
func roomStartPlay(args []interface{}) {
	req := args[0].(*msg.RoomStartPlay)
	log.Debug("房间开始游戏消息-------------------房间ID: %v", req.RoomID)
	// Manager.StartPlay(req.RoomID)
}

// 加入准备房间
func roomAdd(args []interface{}) {
	req := args[0].(*pmsg.AddRequest)
	a := args[1].(gate.Agent)
	log.Debug("房间加入消息----------------房间ID: %v------------------%v", req.RoomID, req.Member.Uuid)
	_, code, err := manager.AddMember(int(req.RoomID), req.Member)
	if err != nil {
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
	manager.LeaveMember(int(req.RoomID), req.Member)
}

// 答题
func roomAnswer(args []interface{}) {
	req := args[0].(*pmsg.Answer)
	log.Debug("收到答题消息------------------------%v", req)
	manager.AnswerQuestion(req)
}

// 房间解散
func roomOver(args []interface{}) {

}

// 房间匹配
func roomMatch(args []interface{}) {
	req := args[0].(*pmsg.MatchRoomRequest)
	log.Debug("房间匹配消息----------------------------------%v", req.RoomID)
	manager.Matching(int(req.RoomID))
}

// 个人匹配
func roomMatchMember(args []interface{}) {
	req := args[0].(*pmsg.MatchMemberRequest)
	log.Debug("房间个人消息----------------------------------%v", req.Member.Uuid)
	room := manager.Create()
	room.AddMember(req.Member)
	room.Matching()
	room.Send(remotemsg.ROOMMATCH, nil)
}

func roomMemberReady(args []interface{}) {
	req := args[0].(*pmsg.MemberReadyRequest)
	log.Debug("房间个人准备消息----------------------------------%v", req.Uuid)
	room, err := manager.MemberReady(int(req.RoomID), req.Uuid)
	if err == nil {
		room.Send(remotemsg.ROOMMEMBERSTATE, nil)
	}
}

func roomMemberReadyCancel(args []interface{}) {
	req := args[0].(*pmsg.MemberReadyCancelRequest)
	log.Debug("房间个人取消准备消息----------------------------------%v", req.Uuid)
	room, err := manager.MemberReadyCancel(int(req.RoomID), req.Uuid)
	if err == nil {
		room.Send(remotemsg.ROOMMEMBERSTATE, nil)
	}
}

func roomMemberLevelChange(args []interface{}) {
	req := args[0].(*pmsg.MemberLevelChange)
	log.Debug("房间个人等级改变消息----------------------------------%v", req.Uuid)
	manager.MemberLevelChange(int(req.RoomID), req.Level, req.Uuid)
}

// 房间匹配取消
func roomMatchCanel(args []interface{}) {
	req := args[0].(*pmsg.MatchRoomCancelRequest)
	bo := manager.MatchingCancel(int(req.RoomID))
	a := args[1].(gate.Agent)
	a.WriteMsg(&pmsg.MatchRoomCancelReply{State: bo})
}

// 个人匹配取消
func roomMatchMemberCanel(args []interface{}) {
	req := args[0].(*pmsg.LeaveRequest)
	manager.MatchingCancel(int(req.RoomID))
	a := args[1].(gate.Agent)
	a.WriteMsg(&pmsg.MatchMemberCancelReply{State: true})
}

// 成员复活
func roomMatchMemberRelive(args []interface{}) {
	req := args[0].(*pmsg.MemberReliveRequest)
	log.Debug("复活消息----------------: %v", req.Uuid)
	manager.Relive(int(req.RoomID), req.Uuid)
	// if err == nil {
	// 	room.SendRelive(req)
	// }
}
