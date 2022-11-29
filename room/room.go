package main

import (
	"catLog"
	"config"
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

func (s *MsgRoom) Over(ctx context.Context, req *msg.OverRequest) (*msg.RoomChangeState, error) {
	room.Manager.OverRoomChan <- int(req.RoomID)
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Leave(ctx context.Context, req *msg.LeaveRequest) (*msg.RoomChangeState, error) {
	room.Manager.LeaveChan <- &room.ChangeRoom{
		RoomID: int(req.RoomID),
		Member: req.Member,
	}
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Add(ctx context.Context, req *msg.AddRequest) (*msg.RoomChangeState, error) {
	// 获得要加入的房间
	room.Manager.AddFriendChan <- &room.ChangeRoom{
		RoomID: int(req.RoomID),
		Member: req.Member,
	}
	return &msg.RoomChangeState{State: 1}, nil
}

func (s *MsgRoom) Create(req *msg.CreateRoomRequest, srv msg.Room_CreateServer) error {
	catLog.Log("******************房间创建请求******************")
	msgMember := req.Member
	subRoom := room.Manager.CreateRoom(msgMember)

	var send = func(msgID int32, m *msg.Member) {
		srv.Send(&msg.CreateRoomReply{
			RoomID:         int32(subRoom.ID),
			PrepareMembers: subRoom.PlayingMembers,
			PlayingMembers: subRoom.PlayingMembers,
			Progress:       int32(subRoom.Cur),
			ChangeMemeber:  m,
			MsgID:          msgID,
		})
	}

loop:
	for {
		select {
		case <-subRoom.OverChan: // 房主解散房间
			catLog.Log("房间结束计时")
			break loop
		case <-subRoom.StartPlayChan: // 游戏开始通知
			send(remotemsg.ROOMSTARTPLAY, nil)
		case <-subRoom.EndPlayChan: // 游戏结束通知
			send(remotemsg.ROOMOVER, nil)
		case changeMember := <-subRoom.AddMemberChan: // 加人通知
			send(remotemsg.ROOMADD, changeMember)
		case changeMember := <-subRoom.LeaveMemberChan: // 离开人通知
			send(remotemsg.ROOMLEAVE, changeMember)
		case changeMember := <-subRoom.PrepareChan: // 准备通知
			send(remotemsg.ROOMPREPARE, changeMember)
		case changeMember := <-subRoom.PrepareCancelChan: // 取消准备通知
			send(remotemsg.ROOMPREPARECANCEL, changeMember)
		case changeMember := <-subRoom.ChangeMasterChan: // 房主变更通知
			send(remotemsg.ROOMCHANGEMASTER, changeMember)
		case <-subRoom.AnswerChan: // 回答问题通知
			send(remotemsg.ROOMANSWEREND, nil)
		case <-subRoom.TimeChan:
			send(remotemsg.ROOMTIME, nil)
		}
	}
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
		close(room.Manager.Done) // 启动失败，关闭通道
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
