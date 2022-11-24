package handler

import (
	"catLog"
	"context"
	"fmt"
	"proto/msg"
	"remotemsg"
	"server/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Storage struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.StorageClient
}

func NewStorage() *Storage {
	s := Storage{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

var StorageInstance *Storage = NewStorage()

func (s *Storage) Run(port string) {
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

type StorageUpdateRequest struct {
	Uid   string `json:"uid"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StoragePullRequest struct {
	Uid string `json:"uid"`
	Key string `json:"key"`
}

func updateStorage(data client.Msg) {
	sr := &StorageUpdateRequest{}
	data.Val.ParseData(sr)
	catLog.Log("更新存储信息", sr.Uid, sr.Key, sr.Value)
	r, err := StorageInstance.innerClient.UpdateStorage(context.Background(), &msg.UpdateStorageRequest{
		Uuid:  StorageInstance.cat.Client.Uuid,
		Uid:   sr.Uid,
		Key:   sr.Key,
		Value: sr.Value,
	})
	if err != nil {
		return
	}
	data.Client.Write(remotemsg.STORAGEUPDATE, r.State)
}

func pullStorage(data client.Msg) {
	sr := &StoragePullRequest{}
	data.Val.ParseData(sr)
	r, err := StorageInstance.innerClient.PullStorage(context.Background(), &msg.PullStorageRequest{
		Uuid: StorageInstance.cat.Client.Uuid,
		Uid:  sr.Uid,
		Key:  sr.Key,
	})
	if err != nil {
		catLog.Log("获取存储内容失败")
		return
	}
	catLog.Log("拉取存储信息", r.Value)

	data.Client.Write(remotemsg.STORAGE, r.Value)
}
