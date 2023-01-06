package internal

import (
	pmsg "proto/msg"
)

type MsgRoom struct {
	pmsg.RoomServer
}

// func (services *MsgRoom) OrderList(params *pmsg.RoomServerConnectRequest, stream pmsg.RoomInfoReply) error {
// 	for {
// 		select{

// 		}
// 	}
// 	return nil
// }
