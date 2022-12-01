package handler

import (
	"catLog"
	"context"
	"fmt"
	"proto/msg"
	"remotemsg"
	"server"
	"server/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Storage struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.StorageClient
	S           *server.Server
}

func NewStorage() *Storage {
	s := Storage{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

var StorageInstance *Storage = NewStorage()

func (s *Storage) Run(port string, ss *server.Server) {
	s.S = ss
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect server failed,", err)
	}
	s.Conn = conn
	s.innerClient = msg.NewStorageClient(conn)

	// 注册消息
	s.cat.Register(remotemsg.STORAGEUPDATE, updateStorage)
	s.cat.Register(remotemsg.STORAGE, pullStorage)

	go func() {
		for {
			select {
			case <-s.cat.GetDone():
				return
			case uuid := <-s.cat.GetOfflineChan():
				state, err := s.innerClient.OffLineStorage(context.Background(), &msg.OffLineStorageRequest{
					Uuid: uuid,
				})
				if err != nil {
					return
				}
				catLog.Log("离线成功_", state.State)
			}

		}
	}()
}

func updateStorage(data client.Msg) {
	sr := &msg.UpdateStorageRequest{}
	data.Val.ParseData(sr)
	r, err := StorageInstance.innerClient.UpdateStorage(context.Background(), sr)
	if err != nil {
		return
	}
	data.Client.MsgChan <- &client.BackMsg{
		MsgID: remotemsg.STORAGEUPDATE,
		Val:   r.State,
	}
}

func pullStorage(data client.Msg) {
	sr := &msg.PullStorageRequest{}
	data.Val.ParseData(sr)
	r, err := StorageInstance.innerClient.PullStorage(context.Background(), sr)
	if err != nil {
		catLog.Log("获取存储内容失败")
		return
	}
	catLog.Log("拉取存储信息", r.Value)
	data.Client.MsgChan <- &client.BackMsg{
		MsgID: remotemsg.STORAGE,
		Val:   r.Value,
	}
}
