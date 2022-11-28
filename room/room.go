package main

import (
	"catLog"
	"context"
	"log"
	"net"
	"proto/msg"
	"remotemsg"
	"room/room"

	"google.golang.org/grpc"
)

type MsgRoom struct {
	msg.RoomServer
}

func (s *MsgRoom) Over(context.Context, *msg.OverRequest) (*msg.RoomChangeState, error) {

	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Leave(ctx context.Context, req *msg.LeaveRequest) (*msg.RoomChangeState, error) {
	err := room.Manager.LeavePrepareMemeber(int(req.RoomID), req.Member)
	if err != nil {
		return nil, err
	}
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Add(ctx context.Context, req *msg.AddRequest) (*msg.RoomChangeState, error) {
	// 获得要加入的房间
	err := room.Manager.AddFriendMember(int(req.RoomID), req.Member)
	if err != nil {
		return nil, err
	}
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Create(req *msg.CreateRoomRequest, srv msg.Room_CreateServer) error {
	catLog.Log("******************房间创建请求******************")
	msgMember := req.Member
	// 创建房主
	member := room.NewMember(
		msgMember.Uuid,
		msgMember.Uid,
		msgMember.Nickname,
		msgMember.Icon,
		true,
		false,
	)
	subRoom := room.Manager.CreateRoom(member)

loop:
	for {
		select {
		case <-subRoom.OverChan: // 房主解散房间
			catLog.Log("房间结束计时")
			break loop
		case <-subRoom.StartPlayChan: // 游戏开始通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				MsgID:          remotemsg.ROOMSTARTPLAY,
			})

		case <-subRoom.EndPlayChan: // 游戏结束通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
			})
		case changeMember := <-subRoom.AddMemberChan: // 加人通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				ChangeMemeber:  changeMember,
			})
		case changeMember := <-subRoom.LeaveMemberChan: // 离开人通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				ChangeMemeber:  changeMember,
			})
		case changeMember := <-subRoom.PrepareChan: // 准备通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				ChangeMemeber:  changeMember,
				MsgID:          remotemsg.ROOMPREPARE,
			})
		case changeMember := <-subRoom.PrepareCancelChan: // 取消准备通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				ChangeMemeber:  changeMember,
				MsgID:          remotemsg.ROOMPREPARECANCEL,
			})
		case changeMember := <-subRoom.ChangeMasterChan: // 房主变更通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				ChangeMemeber:  changeMember,
				MsgID:          remotemsg.ROOMCHANGEMASTER,
			})
		case <-subRoom.AnswerChan: // 回答问题通知
			srv.Send(&msg.CreateRoomReply{
				RoomID:         int32(subRoom.ID),
				PrepareMembers: subRoom.PlayingMembers,
				PlayingMembers: subRoom.PlayingMembers,
				Progress:       int32(subRoom.Cur),
				MsgID:          remotemsg.ROOMANSWEREND,
			})
		}

	}
	return nil
}

const (
	// Address 监听地址
	Address string = ":50056"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	// 默认单次接收最大消息长度为`1024*1024*4`bytes(4M)，单次发送消息最大长度为`math.MaxInt32`bytes
	// grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*4), grpc.MaxSendMsgSize(math.MaxInt32))
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	msg.RegisterRoomServer(grpcServer, &MsgRoom{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
	// // 模拟数据
	// var f = func() {
	// 	// 加入成员
	// 	for i := 0; i < config.ReadRoom().MaxMember*2; i++ {
	// 		uuid := "player" + fmt.Sprint(i)
	// 		uid := "ppp" + fmt.Sprint(i)
	// 		room.MatchManager.AddMatchMember(room.NewMember(uuid, uid, "小猫", "icon", true, false))
	// 	}
	// }
	// f()
	// done := make(chan interface{})
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case <-time.After(time.Second * time.Duration(15)):
	// 		// catLog.Log("挂起******************************")
	// 		f()
	// 	}
	// }

}
