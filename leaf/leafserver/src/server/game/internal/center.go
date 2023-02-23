package internal

import (
	"config"
	"context"
	pmsg "proto/msg"
	"time"

	"github.com/name5566/leaf/log"
	"google.golang.org/grpc"
)

var (
	streamClient  pmsg.CenterClient
	CenterConnect bool
	ConnectCount  int
)

// 中心服初始化
func CenterInit() {
	ConnectCount = 0
	log.Release("连接中心服路径: %v", config.CenterUrl)
	connect, err := grpc.Dial(config.CenterUrl, grpc.WithInsecure())
	if err != nil {
		log.Error("%v", err)
	}
	streamClient = pmsg.NewCenterClient(connect)
	skeleton.Go(func() {
		HearbeatCenter()
	}, func() {})
}

// 中心服心跳
func HearbeatCenter() {
	//调用服务端RouteList方法，获流
	stream, err := streamClient.Heartbeat(context.Background())
	if err != nil {
		log.Error("Upload list err: %v", err)
		if ConnectCount <= 0 {
			time.Sleep(time.Second * time.Duration(5))
			if ConnectCount <= 0 {
				HearbeatCenter()
			}
		}
		return
	}
	CenterConnect = true
	ConnectCount++
	for {
		err := stream.Send(&pmsg.CenterPing{Url: ServerConf.SelfUrl + config.Port, Count: int32(len(Users))})
		if err != nil {
			log.Error("stream request err: %v", err)
			break
		}
		time.Sleep(time.Second * time.Duration(5))
	}
	//关闭流并获取返回的消息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Error("%v", err)
	}
	// log.Println(res)
	log.Debug("关闭中心服流: %v", res)
	CenterConnect = false
	for {
		time.Sleep(time.Second * time.Duration(5))
		if !CenterConnect {
			HearbeatCenter()
		} else {
			log.Debug("退出连接监听**********************")
			break
		}
	}
}
