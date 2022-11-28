package handler

import (
	"catLog"
	"context"
	"fmt"
	"io"
	"log"
	"proto/msg"
	"remotemsg"
	"room/room"
	"server/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Room struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.RoomClient
}

var RoomInstance *Room = NewRoom()

func NewRoom() *Room {
	s := Room{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

func CreateRoom() {
	var member = room.NewMember("dafafa", "cat", "小猫", "icon", true, false)
	// 创建发送结构体
	req := msg.CreateRoomRequest{
		Member: member,
	}

	// 调用我们的服务(ListValue方法)
	stream, err := RoomInstance.innerClient.Create(context.Background(), &req)
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
		}
	}()
}

func (s *Room) Run(port string) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect server failed,", err)
	}
	s.Conn = conn
	s.innerClient = msg.NewRoomClient(conn)
	CreateRoom()

	// 注册消息
	s.cat.Register(remotemsg.ROOMCREATE, roomCreate)
	s.cat.Register(remotemsg.ROOMSTARTPLAY, roomStartPlay)
	s.cat.Register(remotemsg.ROOMADD, roomAdd)
	s.cat.Register(remotemsg.ROOMLEAVE, roomLeave)
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

type RoomBaseRequest struct {
	Uid      string // 用户Uid
	Nickname string // 用户昵称
	Icon     string // 用户头像
}

func roomCreate(data client.Msg) {
	u := &RoomBaseRequest{}
	data.Val.ParseData(u)
	CreateRoom()
}

func roomStartPlay(data client.Msg) {
	u := &RoomBaseRequest{}
	data.Val.ParseData(u)
}

func roomAdd(data client.Msg) {
	u := &RoomBaseRequest{}
	data.Val.ParseData(u)
}

func roomLeave(data client.Msg) {
	u := &RoomBaseRequest{}
	data.Val.ParseData(u)
}

func roomAnswer(data client.Msg) {
	u := &RoomBaseRequest{}
	data.Val.ParseData(u)
}
