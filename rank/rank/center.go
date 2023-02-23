package rank

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"proto/msg"
	pmsg "proto/msg"
	"sort"
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
	Servers     []*GameServer
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
		fmt.Printf("在线人数: %v|%v\n", req.Url, req.Count)
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
			Servers = append(Servers, obj)
		}
	}
}

func SortServers() {
	sort.SliceStable((Servers), func(i, j int) bool {
		return (Servers)[i].Count < (Servers)[j].Count
	})
}

func CenterInit() {
	println("启动rpc监听")
	GameServers = make(map[string]*GameServer)
	Servers = []*GameServer{}
	go func() {
		for {
			lock.Lock()
			for _, v := range GameServers {
				v.HeatbeatTime++
				if v.HeatbeatTime >= 3 {
					fmt.Printf("删除:----%v\n", v.Url)
					delete(GameServers, v.Url) // 删除
					Servers = deleteServer(Servers, v)
				}
			}
			SortServers() // 每隔5秒排序
			lock.Unlock()
			time.Sleep(time.Second * time.Duration(5))
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
	err = grpcServer.Serve(listener)
	if err != nil {
		// log.Fatalf("failed to serve: %v", err)

	}
}

func RoomCreate(c echo.Context) error {
	lock.RLock()
	var url = ""
	if len(Servers) > 0 {
		url = Servers[0].Url
		println("下发路径", url, Servers[0].Count)
	}
	lock.RUnlock()
	// url := Conf.Urls[curServer]
	// curServer++
	// if curServer >= serversMax {
	// 	curServer = 0
	// }
	return c.JSON(http.StatusOK, &pmsg.RoomPreAddReply{
		Url: url,
	})
}

func deleteServer(a []*GameServer, elem *GameServer) []*GameServer {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
