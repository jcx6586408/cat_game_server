package main

import (
	"catLog"
	"config"
	"context"
	"log"
	"net"
	"proto/msg"
	"room/room"

	"google.golang.org/grpc"
)

type MsgRoom struct {
	msg.RoomServer
}

// 离线
func (s *MsgRoom) Offline(ctx context.Context, req *msg.OfflineRequest) (*msg.RoomChangeState, error) {
	room.Manager.OfflineChan <- req.Uuid
	return &msg.RoomChangeState{State: 1}, nil
}

// 个人匹配取消
func (s *MsgRoom) MatchMemberCancel(ctx context.Context, req *msg.LeaveRequest) (*msg.RoomChangeState, error) {
	room.Manager.MatchMemberCancelChan <- req
	return &msg.RoomChangeState{State: 1}, nil
}

// 取消匹配房间
func (s *MsgRoom) MatchRoomCancel(ctx context.Context, req *msg.MatchRoomRequest) (*msg.RoomChangeState, error) {
	room.Manager.MatchRoomCancelChan <- int(req.RoomID)
	return &msg.RoomChangeState{State: 1}, nil
}

// 匹配房间
func (s *MsgRoom) MatchRoom(ctx context.Context, req *msg.MatchRoomRequest) (*msg.RoomChangeState, error) {
	room.Manager.MathRoomChan <- int(req.RoomID)
	return &msg.RoomChangeState{State: 1}, nil
}

// 匹配单人
func (s *MsgRoom) MatchMember(ctx context.Context, req *msg.MatchMemberRequest) (*msg.RoomChangeState, error) {
	room.Manager.MathMemberChan <- req.Member
	return &msg.RoomChangeState{State: 1}, nil
}

// 回答问题
func (s *MsgRoom) AnswerQuestion(ctx context.Context, req *msg.Answer) (*msg.RoomChangeState, error) {
	room.Manager.AnswerChan <- req
	return &msg.RoomChangeState{State: 1}, nil
}

// 结束房间
func (s *MsgRoom) Over(ctx context.Context, req *msg.OverRequest) (*msg.RoomChangeState, error) {
	room.Manager.OverRoomChan <- int(req.RoomID)
	return &msg.RoomChangeState{State: 1}, nil
}

// 离开房间
func (s *MsgRoom) Leave(ctx context.Context, req *msg.LeaveRequest) (*msg.RoomChangeState, error) {
	room.Manager.LeaveChan <- &room.ChangeRoom{
		RoomID: int(req.RoomID),
		Member: req.Member,
	}
	return &msg.RoomChangeState{State: 1}, nil
}

// 加入房间
func (s *MsgRoom) Add(ctx context.Context, req *msg.AddRequest) (*msg.RoomChangeState, error) {
	catLog.Log(req.Member.Uuid, "请求加入房间", req.Member.IsMaster)
	// 获得要加入的房间
	room.Manager.AddFriendChan <- &room.ChangeRoom{
		RoomID: int(req.RoomID),
		Member: req.Member,
	}
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Create(ctx context.Context, req *msg.CreateRoomRequest) (*msg.RoomChangeState, error) {
	msgMember := req.Member
	room.Manager.CreateChan <- msgMember
	subRoom := <-room.Manager.CreateCompleteChan
	catLog.Log("*************创建房间成功**************", subRoom.ID)
	return &msg.RoomChangeState{State: 1}, nil
}

func send(subRoom *room.Room, srv msg.Room_ConnectServer, msgID int32, m *msg.Member) {
	srv.Send(&msg.CreateRoomReply{
		RoomID:         int32(subRoom.ID),
		PrepareMembers: subRoom.PrepareMembers,
		PlayingMembers: subRoom.PlayingMembers,
		Progress:       int32(subRoom.Cur),
		TotolQuestion:  int32(subRoom.QuestionCount),
		CurQuestion:    int32(subRoom.GetProgress()),
		ChangeMemeber:  m,
		Question:       subRoom.GetQuestion(),
		MsgID:          msgID,
		ToTalTime:      int32(subRoom.GetPlayTime()),
	})
}

// 创建房间
func (s *MsgRoom) Connect(req *msg.RoomServerConnectRequest, srv msg.Room_ConnectServer) error {

	// catLog.Log("房间创建通知", subRoom.ID)
	// send(remotemsg.ROOMCREATE, nil) // 首次创建通知

	// // loop:
	// for {
	// 	select {
	// 	// case <-
	// 	case <-subRoom.OverChan: // 房主解散房间
	// 		catLog.Log("房主解散房间")
	// 		send(remotemsg.ROOMOVER, nil)
	// 	case <-subRoom.StartPlayChan: // 游戏开始通知
	// 		catLog.Log("游戏开始通知*********************")
	// 		send(remotemsg.ROOMSTARTPLAY, nil)
	// 	case <-subRoom.EndPlayChan: // 游戏结束通知
	// 		send(remotemsg.ROOMOVER, nil)
	// 	case changeMember := <-subRoom.AddMemberChan: // 加人通知
	// 		send(remotemsg.ROOMADD, changeMember)
	// 	case changeMember := <-subRoom.LeaveMemberChan: // 离开人通知
	// 		catLog.Log("成员离开通知", changeMember.Uuid)
	// 		send(remotemsg.ROOMLEAVE, changeMember)
	// 	case changeMember := <-subRoom.PrepareChan: // 准备通知
	// 		send(remotemsg.ROOMPREPARE, changeMember)
	// 	case changeMember := <-subRoom.PrepareCancelChan: // 取消准备通知
	// 		send(remotemsg.ROOMPREPARECANCEL, changeMember)
	// 	case changeMember := <-subRoom.ChangeMasterChan: // 房主变更通知
	// 		send(remotemsg.ROOMCHANGEMASTER, changeMember)
	// 	case <-subRoom.AnswerChan: // 回答问题通知
	// 		send(remotemsg.ROOMANSWEREND, nil)
	// 	case changeMember := <-subRoom.MemberAnswerChan: // 回答问题通知
	// 		send(remotemsg.ROOMANSWER, changeMember)
	// 	case <-subRoom.TimeChan: // 计时通知
	// 		send(remotemsg.ROOMTIME, nil)
	// 	case changeMember := <-subRoom.OfflineChan: // 离线通知
	// 		send(remotemsg.ROOMTIME, changeMember)
	// 	}
	// }
	return nil
}

var (
	// Address 默认监听地址
	Address string = ":50056"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	room.Manager.Run() // 开启监听
	var conf = config.Read()
	Address = conf.Room.Port // 配置端口
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		close(room.Manager.Done) // 启动失败，关闭通道
		log.Fatalf("net.Listen err: %v", err)
		return
	}
	catLog.Log(Address + " net.Listing...")
	// 新建gRPC服务器实例
	// 默认单次接收最大消息长度为`1024*1024*4`bytes(4M)，单次发送消息最大长度为`math.MaxInt32`bytes
	// grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*4), grpc.MaxSendMsgSize(math.MaxInt32))
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	msg.RegisterRoomServer(grpcServer, &MsgRoom{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		close(room.Manager.Done) // 启动失败，关闭通道
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
