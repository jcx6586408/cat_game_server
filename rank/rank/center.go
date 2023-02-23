package rank

import (
	"io"
	"log"
	"net"
	"net/http"
	"proto/msg"
	pmsg "proto/msg"
	"storage/redis"
	"time"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Serve struct {
	msg.CenterServer
}

// 游戏服
type GameServer struct {
	Url          string // 游戏服路径
	Count        int    // 游戏服人数
	HeatbeatTime int
}

var (
	GameServers map[string]*GameServer // 游戏服
	redisKey    string                 = "gameservers"
	Port        string                 // 端口监听
)

func (s *Serve) Heartbeat(srv msg.Center_HeartbeatServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			println("读取结束", err)
			return nil
		}
		if err != nil {
			println("读取错误", err)
			return err
		}
		// fmt.Printf("在线人数: %v|%v\n", req.Url, req.Count)
		redis.AddGameServers(redisKey, req.Url, float64(req.Count))
		obj, ok := GameServers[req.Url]
		if ok {
			obj.Count = int(req.Count)
			obj.HeatbeatTime = 0
		} else {
			obj = &GameServer{
				Url:          req.Url,
				Count:        int(req.Count),
				HeatbeatTime: 0}
			GameServers[req.Url] = obj
		}
	}
}

func CenterInit() {
	println("启动rpc监听")
	GameServers = make(map[string]*GameServer)
	go func() {
		for {
			for _, v := range GameServers {
				v.HeatbeatTime++
				if v.HeatbeatTime >= 3 {
					redis.DeleGameServer(redisKey, v.Url)
					delete(GameServers, v.Url) // 删除
				}
			}
			time.Sleep(time.Second * time.Duration(5))
			// fmt.Printf("服务器: %v\n", redis.GetTopGameServers(redisKey)[0].Member)
		}
	}()
	// 创建 Tcp 连接
	listener, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	// 创建gRPC服务
	grpcServer := grpc.NewServer()

	// 注册游戏服务
	msg.RegisterCenterServer(grpcServer, &Serve{})

	// 在 gRPC 服务上注册反射服务
	reflection.Register(grpcServer)

	println("********中心服启动成功********0")
	// 连接redis
	redis.ConnectReids()
	err = grpcServer.Serve(listener)
	if err != nil {
		// log.Fatalf("failed to serve: %v", err)

	}
}

func RoomCreate(c echo.Context) error {
	// url := Conf.Urls[curServer]
	// curServer++
	// if curServer >= serversMax {
	// 	curServer = 0
	// }
	return c.JSON(http.StatusOK, &pmsg.RoomPreAddReply{
		Url: redis.GetTopGameServers(redisKey)[0].Member.(string),
	})
}
