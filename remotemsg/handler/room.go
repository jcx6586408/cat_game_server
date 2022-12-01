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
	catLog.Log("开始房间监听", stream)
	go func() {
		for {
			//Recv() 方法接收服务端消息，默认每次Recv()最大消息长度为`1024*1024*4`bytes(4M)
			res, err := stream.Recv()
			// 判断消息流是否已经结束
			if err == io.EOF {
				catLog.Log("房间结束")
				break
			}
			if err != nil {
				log.Fatalf("ListStr get stream err: %v", err)

			}
			// 打印返回值
			catLog.Log("返回房间消息——", res.RoomID)
			// 通知准备的
			for _, v := range res.PrepareMembers {
				catLog.Log("返回准备成员消息----", v.Uid)
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
				catLog.Log("返回正在玩成员消息----", v.Uid)
				c, ok := RoomInstance.S.GetClient(v.Uuid)
				if ok {
					switch res.MsgID {
					case remotemsg.ROOMANSWER:
						c.MsgChan <- &client.BackMsg{
							MsgID: int(res.MsgID),
							Val:   res.ChangeMemeber.Answer,
						}
					default:
						c.MsgChan <- &client.BackMsg{
							MsgID: int(res.MsgID),
							Val:   res,
						}

					}
				}
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
	s.cat.Register(remotemsg.ROOMCREATE, roomCreate)    // 房间创建
	s.cat.Register(remotemsg.ROOMADD, roomAdd)          // 房间加入玩家
	s.cat.Register(remotemsg.ROOMLEAVE, roomLeave)      // 离开房间
	s.cat.Register(remotemsg.ROOMANSWEREND, roomAnswer) // 答题注册

	// 解散房间注册
	s.cat.Register(remotemsg.ROOMOVER, roomOver)

	// 单人匹配与取消注册
	s.cat.Register(remotemsg.ROOMMATCHROOMCANCEL, roomMatchCanel)
	s.cat.Register(remotemsg.ROOMMATCHROOM, roomMatch)

	// 房间准备与取消匹配注册
	s.cat.Register(remotemsg.ROOMMATCH, roomMatchMember)
	s.cat.Register(remotemsg.ROOMMATCHMEMBERCANCEL, roomMatchMemberCanel)

	go func() {
		for {
			select {
			case <-s.cat.GetDone():
				return
			case uuid := <-s.cat.GetOfflineChan():
				catLog.Log("玩家离线, uuid:", uuid)
				u := &msg.OfflineRequest{Uuid: uuid}
				RoomInstance.innerClient.Offline(context.Background(), u)
			}

		}
	}()
}

// 房间创建
func roomCreate(data client.Msg) {
	u := &msg.CreateRoomRequest{}
	data.Val.ParseData(u)
	catLog.Log("房间创建监听============", u.Member)
	CreateRoom(u)
}

// 加入准备房间
func roomAdd(data client.Msg) {
	u := &msg.AddRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.Add(context.Background(), u)
}

// 离开准备房间
func roomLeave(data client.Msg) {
	u := &msg.LeaveRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.Leave(context.Background(), u)
}

// 答题
func roomAnswer(data client.Msg) {
	u := &msg.Answer{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.AnswerQuestion(context.Background(), u)
}

// 房间解散
func roomOver(data client.Msg) {
	u := &msg.OverRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.Over(context.Background(), u)
}

// 房间匹配
func roomMatch(data client.Msg) {
	u := &msg.MatchRoomRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.MatchRoom(context.Background(), u)
}

// 个人匹配
func roomMatchMember(data client.Msg) {
	u := &msg.MatchMemberRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.MatchMember(context.Background(), u)
}

// 房间匹配取消
func roomMatchCanel(data client.Msg) {
	u := &msg.MatchRoomRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.MatchRoomCancel(context.Background(), u)
}

// 个人匹配取消
func roomMatchMemberCanel(data client.Msg) {
	u := &msg.LeaveRequest{}
	data.Val.ParseData(u)
	RoomInstance.innerClient.MatchMemberCancel(context.Background(), u)
}
