package main

import (
	"catLog"
	"config"
	"remotemsg/handler"
	"server"
	"server/client"
)

type Gate struct {
	server.Server // 远程服务
}

func main() {
	s := server.New()
	s.Name = "gate"
	conf := config.Read()
	s.ID = conf.Gate.ID
	s.Port = conf.Gate.Port
	s.HttpsPort = conf.Gate.HttpsPort
	s.CertFile = conf.Crt.CertFile
	s.KeyFile = conf.Crt.KeyFile
	s.ListenerType = conf.Gate.ListenerType
	s.MsgHandler = MsgHandler
	s.UserHandler = UserHandler

	handler.RankInstance.Run(conf.Rank.Port)          // 排行榜注册
	handler.StorageInstance.Run(conf.Storage.Port, s) // 存储注册
	handler.RoomInstance.Run(conf.Room.Port, s)       // 房间注册
	s.Run()                                           // 运行中心服
}

func UserHandler(c *client.Client) {
	handler.RegisterOffline(c) // 设置客户端监听离线和登录信息
}

func MsgHandler(data []byte, c *client.Client) {
	msg := &client.Msg{}
	msg.Client = c
	subMsg := client.NewSubMsg(string(data))
	msg.Val = subMsg
	// 获取处理器
	handler, ok := client.GetHanlder(subMsg.ID)
	if ok {
		Chan := handler.Chan
		Chan <- *msg
	} else {
		catLog.Warn("未注册消息", subMsg.ID)
	}
}
