package client

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Uuid string
	Ws   *websocket.Conn
	R    *http.Request
	C    chan interface{}
	Lock sync.Mutex
}

var clients = make(map[string]*Client)

func New(uuid string, ws *websocket.Conn, r *http.Request) *Client {
	c := &Client{}
	c.Uuid = uuid
	c.Ws = ws
	c.R = r
	c.C = make(chan interface{})
	// 加入客户端
	AddClient(c)
	return c
}

// GetClient 获取客户端
func GetClient(uid string) (c *Client, ok bool) {
	c, ok = clients[uid]
	return c, ok
}

// AddClient 添加客户端链接
func AddClient(client *Client) {
	clients[client.Uuid] = client
}

// Dele 删除链接
func Dele(uuid string) {
	client, ok := clients[uuid]
	if ok {
		close(client.C) // 通知当前连接已经关闭
	}
	delete(clients, uuid)
}
