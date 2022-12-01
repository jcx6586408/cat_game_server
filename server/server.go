package server

import (
	"catLog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"remotemsg"
	"server/client"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server struct {
	ID             int                                // 服务器ID
	Name           string                             // 服务器名字
	Host           string                             // 监听域名
	ListenerType   []string                           // 监听类型
	Port           string                             // 监听http端口
	HttpsPort      string                             // 监听https端口
	CertFile       string                             // 证书crt
	KeyFile        string                             // 证书key
	upgrader       websocket.Upgrader                 // websocket升级类
	UserHandler    func(c *client.Client)             // 用户处理器
	MsgHandler     func(msg []byte, c *client.Client) // 消息处理器
	Clients        map[string]*client.Client          // 接入的客户端
	ConnectChan    chan *client.Client                // 客户端接入管道
	DisConnectChan chan *client.Client                // 客户端断开连接管道
	Done           chan interface{}                   // 关闭服务
}

func New() *Server {
	s := &Server{}
	s.Done = make(chan interface{})
	s.ConnectChan = make(chan *client.Client)
	s.DisConnectChan = make(chan *client.Client)
	s.Clients = make(map[string]*client.Client)
	return s
}

// Dele 删除链接
func (s *Server) Dele(c *client.Client) {
	s.DisConnectChan <- c
}

// 获取客户端
func (s *Server) GetClient(uuid string) (*client.Client, bool) {
	client, ok := s.Clients[uuid]
	return client, ok
}

func (s *Server) Run() {
	http.HandleFunc("/ws", s.wsEndpoint)
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// 启动客户端接入处理
	go func() {
		for {
			select {
			case <-s.Done:
				return
			case c := <-s.ConnectChan:
				s.Clients[c.Uuid] = c
				catLog.Log("接入连接", c.Uuid)
				s.UserHandler(c)
			case c := <-s.DisConnectChan:
				uuid := c.Uuid
				catLog.Log("断开连接", uuid)
				c.Ws.Close()
				close(c.C) // 通知当前连接已经关闭
				_, ok := s.Clients[uuid]
				if ok {
					delete(s.Clients, uuid)
				}
			}
		}
	}()

	// 启动监听
	for _, v := range s.ListenerType {
		if v == "http" {
			catLog.Log(s.Name, "启动http", "监听端口", s.Port)
			log.Fatal(http.ListenAndServe(s.Port, nil))
		}
		if v == "https" {
			catLog.Log(s.Name, "启动https", "监听端口", s.Port)
			log.Fatal(http.ListenAndServeTLS(s.HttpsPort, s.CertFile, s.KeyFile, nil))
		}
	}
}

func (s *Server) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// 跨域设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	s.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		catLog.Warn("升级失败")
		return
	}

	var done = make(chan interface{})
	var interrupt = make(chan os.Signal, 1)
	guid := uuid.New()
	c := client.New(guid.String(), ws, r) // 创建客户端
	s.ConnectChan <- c
	// 登录成功通知，并且下发uuid作为客户端消息凭证
	c.Write(remotemsg.LOGIN, c.Uuid)
	// 通知打断
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// 删除客户端
	defer s.Dele(c)

	// 启动消息读取
	go s.reader(done, ws, c)

	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			err := c.Write(remotemsg.HEARTBEAT, []byte("hb"))
			if err != nil {
				catLog.Warn("没有监听到心跳消息, 断开连接:", err)
				return
			}

		case <-interrupt:
			catLog.Warn("接收到打断的消息，关闭连接")
			err := c.WriteClose()
			if err != nil {
				catLog.Warn("Error during closing websocket:", err)
				return
			}

			select {
			case <-done:
				catLog.Warn("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				catLog.Warn("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

func (s *Server) reader(done chan interface{}, conn *websocket.Conn, c *client.Client) {
	defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			catLog.Warn("读取消息失败", err)
			return
		}
		s.MsgHandler(msg, c)
	}
}
