package handler

import (
	"catLog"
	"context"
	"fmt"
	"io"
	"log"
	"proto/msg"
	"remotemsg"
	"server"
	"server/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Room struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.RoomClient
	S           *server.Server
}

var RoomInstance *Room = NewRoom()

func NewRoom() *Room {
	s := Room{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

func CreateRoom(req *msg.CreateRoomRequest) {
	// 调用我们的服务(ListValue方法)
	stream, err := RoomInstance.innerClient.Create(context.Background(), req)
	if err != nil {
		log.Fatalf("Call ListStr err: %v", err)
	}
	go func() {
		for {
			//Recv() 方法接收服务端消息，默认每次Recv()最大消息长度为`1024*1024*4`bytes(4M)
			res, err := stream.Recv()
			// 判断消息流是否已经结束
			if err == io.EOF {
				catLog.Log("房间结束", res.RoomID)
				break
			}
			if err != nil {
				log.Fatalf("ListStr get stream err: %v", err)

			}
			// 打印返回值
			catLog.Log("返回房间消息——", res.RoomID)
			// 通知准备的
			for _, v := range res.PrepareMembers {
				c, ok := RoomInstance.S.GetClient(v.Uuid)
				if ok {
					c.MsgChan <- &client.BackMsg{
						MsgID: int(res.MsgID),
						Val:   res,
					}
				}
			}
			// 通知正在玩的
			for _, v := range res.PlayingMembers {
				c, ok := RoomInstance.S.GetClient(v.Uuid)
				if ok {
					c.MsgChan <- &client.BackMsg{
						MsgID: int(res.MsgID),
						Val:   res,
					}
				}
			}
			switch res.MsgID {
			case remotemsg.ROOMADD:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMLEAVE:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMPREPARE:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMPREPARECANCEL:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMOVER:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMSTARTPLAY:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMENDPLAY:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMANSWEREND:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			case remotemsg.ROOMTIME:
				catLog.Log("玩家加入_", res.ChangeMemeber.Uid)
			}
		}
	}()
}

func (s *Room) Run(port string, ss *server.Server) {
	s.S = ss
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect server failed,", err)
	}
	s.Conn = conn
	s.innerClient = msg.NewRoomClient(conn)

	// 注册消息
	s.cat.Register(remotemsg.ROOMCREATE, roomCreate)          // 房间创建
	s.cat.Register(remotemsg.ROOMMATCHROOM, roomStartMatch)   // 房间开始匹配
	s.cat.Register(remotemsg.ROOMMATCH, roomStartMemberMatch) // 房间个人开始匹配
	s.cat.Register(remotemsg.ROOMADD, roomAdd)                // 房间加入玩家
	s.cat.Register(remotemsg.ROOMLEAVE, roomLeave)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)

	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer)

	go func() {
		for {
			select {
			case <-s.cat.GetDone():
				return
			case uuid := <-s.cat.GetOfflineChan():
				catLog.Log("玩家离线, uuid:", uuid)
				// state, err := s.innerClient.OffLineStorage(context.Background(), &msg.OffLineStorageRequest{
				// 	Uuid: uuid,
				// })
				// if err != nil {
				// 	return
				// }
				// catLog.Log("离线成功_", state.State)

			}

		}
	}()
}

func roomCreate(data client.Msg) {
	u := &msg.CreateRoomRequest{}
	data.Val.ParseData(u)
	CreateRoom(u)
}

// 房间匹配
func roomStartMatch(data client.Msg) {
	u := &msg.MatchRoomRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.MatchRoom(context.Background(), u)
}

// 个人匹配
func roomStartMemberMatch(data client.Msg) {
	// u := &RoomBaseRequest{}
	// data.Val.ParseData(u)
}

// 加入准备房间
func roomAdd(data client.Msg) {
	u := &msg.AddRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.Add(context.Background(), u)
}

func roomLeave(data client.Msg) {
	u := &msg.LeaveRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.Leave(context.Background(), u)
}

func roomAnswer(data client.Msg) {
	u := &msg.Answer{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.AnswerQuestion(context.Background(), u)
}
