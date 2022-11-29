package client

import (
	"catLog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Uuid    string
	Ws      *websocket.Conn
	R       *http.Request
	C       chan interface{} // 关闭连接通知通道
	Lock    sync.Mutex
	MsgChan chan *BackMsg // 消息写入
}

var lock sync.RWMutex

var ClientPool = &sync.Pool{
	New: func() interface{} {
		return &Client{}
	},
}

func New(uuid string, ws *websocket.Conn, r *http.Request) *Client {
	c := &Client{}
	c.Uuid = uuid
	c.Ws = ws
	c.R = r
	c.C = make(chan interface{})
	c.MsgChan = make(chan *BackMsg)
	// 监听消息写入
	go func() {
		for {
			select {
			case <-c.C:
				catLog.Log("关闭消息监听")
				return
			case msg := <-c.MsgChan:
				c.Write(msg.MsgID, msg.Val)
			}
		}
	}()
	return c
}
