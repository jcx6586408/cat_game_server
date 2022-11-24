package client

import (
	"catLog"
	"encoding/json"
	"fmt"
)

type SubMsg struct {
	Uuid string `json:"uuid"` // 玩家uuid（服务器下发，客户端存储）
	ID   int    `json:"id"`   // 消息ID
	Data string `json:"data"` // 消息体
}

func NewSubMsg(msg string) *SubMsg {
	s := &SubMsg{}
	json.Unmarshal([]byte(msg), s)
	catLog.Log("解析后数据", s.Data)
	return s
}

func (subMsg *SubMsg) GetData() string {
	return subMsg.Data
}

func (subMsg *SubMsg) ParseData(t interface{}) {
	json.Unmarshal([]byte(subMsg.Data), t)
}

func (subMsg *SubMsg) ToString() string {
	return fmt.Sprint(subMsg.ID) + subMsg.Data
}
