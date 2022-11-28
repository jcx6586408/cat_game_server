package client

import "catLog"

// 消息处理器
type MsgHandler struct {
	MsgID int // 消息ID
	Chan  chan Msg
}

var msgHanlders = make(map[int]*MsgHandler)

// AddHandler 加入消息处理器
func RegisterHandler(handler *MsgHandler) {
	catLog.Log("注册消息ID: ", handler.MsgID)
	lock.Lock()
	msgHanlders[handler.MsgID] = handler
	lock.Unlock()
}

// GetHanlder 获取处理器
func GetHanlder(id int) (*MsgHandler, bool) {
	h, ok := msgHanlders[id]
	return h, ok
}
