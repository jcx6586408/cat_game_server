package internal

import (
	pmsg "proto/msg"
)

var (
	RoomCreateChan chan int                 // 房间创建通道
	RoomAddChan    chan *pmsg.Member        // 房间加入通道
	RoomLeaveChan  chan *pmsg.Member        // 房间离开通道
	RoomChangeChan chan *pmsg.RoomInfoReply // 房间状态变更通道
	Done           chan interface{}         // 服务器关闭
)

type MsgRoom struct {
	pmsg.RoomServer
}

func (services *MsgRoom) OrderList(params *pmsg.RoomServerConnectRequest, stream pmsg.Room_ConnectServer) error {
loop:
	for {
		select {
		case <-Done:
			break loop
		case info := <-RoomChangeChan:
			stream.Send(info)
		}
	}
	return nil
}
