package storage

import (
	"catLog"
	"context"
	"proto/msg"

	"sync"
)

type Storage struct {
	msg.StorageServer
}

func New() *Storage {
	s := Storage{}
	return &s
}

var StorageInstance *Storage = New()

func (s *Storage) LoginStorage(ctx context.Context, r *msg.LoginStorageRequest) (*msg.SuccessStateReply, error) {
	catLog.Log("用户登录", r.Uid)
	return &msg.SuccessStateReply{State: 1}, nil
}

func (s *Storage) OffLineStorage(ctx context.Context, r *msg.OffLineStorageRequest) (*msg.SuccessStateReply, error) {
	catLog.Log("用户离线", r.Uuid)
	s.OnClose(r.Uuid)
	return &msg.SuccessStateReply{State: 1}, nil
}

func (s *Storage) UpdateStorage(ctx context.Context, r *msg.UpdateStorageRequest) (*msg.SuccessStateReply, error) {
	catLog.Log("用户更新存储", r.Uid, r.Nickname, r.Icon, r.Key, r.Value)
	backInfo := NewUserStorage(r.Uid, r.Uuid, r.Nickname, r.Icon)
	backInfo.Forever[r.Key] = r.Value
	catLog.Log("存储信息展示")
	for k, v := range backInfo.Forever {
		catLog.Log(k, v)
	}
	return &msg.SuccessStateReply{State: 1}, nil
}

func (s *Storage) PullStorage(ctx context.Context, r *msg.PullStorageRequest) (*msg.PullStorageReply, error) {
	user := NewUserStorage(r.Uid, r.Uuid, r.Nickname, r.Icon)
	backInfo, ok := user.Forever[r.Key]
	catLog.Log("用户拉取存储", r.Uid, "存储内容", backInfo)
	if ok {
		return &msg.PullStorageReply{Value: backInfo}, nil
	} else {
		return &msg.PullStorageReply{Value: ""}, nil
	}
}

func (s *Storage) OnClose(uuid string) {
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		defer w.Done()
		for _, v := range users {
			if v.Uuid == uuid {
				catLog.Log("玩家", v.Uid, "下线", v.Forever)
				// 保存玩家信息
				c := SaveToDB(v)
				<-c
				delete(users, v.Uid)
				break
			}
		}
	}()
	w.Wait()
	catLog.Log("当前在线玩家数量_", len(users))
}

func (s *Storage) Run() {
	catLog.Log("开始进行数据库连接")
	Connect() // 连接数据库
}

type UserStorageDB struct {
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
	Icon     string `json:"icon"`
	Forever  string `json:"forever"`
}

type UserStorage struct {
	Uuid     string            `json:"uuid"`
	Uid      string            `json:"uid"`
	Nickname string            `json:"nickname"`
	Icon     string            `json:"icon"`
	Forever  map[string]string `json:"forever"`
}

func NewUserStorage(uid string, uuid string, nickname string, icon string) *UserStorage {
	_, ok := users[uid]
	if !ok {
		users[uid] = &UserStorage{
			Uuid:     uuid,
			Uid:      uid,
			Nickname: nickname,
			Icon:     icon,
			Forever:  make(map[string]string),
		}
		// forever
		c := userFromDb(uid, users[uid])
		<-c
		catLog.Log("当前在线玩家数量_", len(users))
	}
	return users[uid]
}

var users = make(map[string]*UserStorage)
