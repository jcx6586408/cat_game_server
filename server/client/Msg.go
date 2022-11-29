package client

import (
	"catLog"
	"encoding/json"

	"github.com/gorilla/websocket"
)

// Msg 消息
type Msg struct {
	Client *Client
	Val    *SubMsg // 上传消息信息
}

type BackMsg struct {
	MsgID int         // 消息ID
	Val   interface{} // 返回消息信息
}

func (c *Client) Write(id int, backInfo interface{}) error {
	var rmsg = NewRemoteMsg(id, 0, backInfo)
	jsons, err := json.Marshal(rmsg)
	if err == nil {
		e := c.WriteMsg(jsons)
		return e
	} else {
		catLog.Warn("解析消息出错", err)
	}
	return err
}

func (c *Client) WriteMsg(data []byte) error {
	c.Lock.Lock()
	err := c.Ws.WriteMessage(websocket.TextMessage, data)
	c.Lock.Unlock()
	return err
}

func (c *Client) WriteClose() error {
	c.Lock.Lock()
	err := c.Ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.Lock.Unlock()
	return err
}

func (c *Client) ReadMsg() (data []byte) {
	_, p, err := c.Ws.ReadMessage()
	if err != nil {
		catLog.Warn("消息读取失败", err)
		return
	}
	return p
}
