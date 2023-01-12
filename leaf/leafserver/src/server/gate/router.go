package gate

import (
	"leafserver/src/server/game"
	"leafserver/src/server/login"
	"leafserver/src/server/msg"
	pmsg "proto/msg"
)

func init() {
	// 排行请求自身信息
	msg.Processor.SetRouter(&msg.RankSelfRequest{}, game.ChanRPC)
	// 排行请求列表
	msg.Processor.SetRouter(&msg.Rank{}, game.ChanRPC)
	// 排行请求更新
	msg.Processor.SetRouter(&msg.RankPull{}, game.ChanRPC)

	// 登录
	msg.Processor.SetRouter(&msg.WXCode{}, game.ChanRPC)

	// msg.ProbufProcessor.SetRouter(&pmsg.RoomInfoGetRequest{}, login.ChanRPC)

	// 房间创建消息
	msg.Processor.SetRouter(&pmsg.RoomPreAddRequest{}, login.ChanRPC) // 房间预创建
	msg.Processor.SetRouter(&pmsg.CreateRoomRequest{}, game.ChanRPC)  // 房间创建
	msg.Processor.SetRouter(&msg.RoomStartPlay{}, game.ChanRPC)       // 房间主动开始游戏
	msg.Processor.SetRouter(&pmsg.AddRequest{}, game.ChanRPC)         // 加入房间注册
	msg.Processor.SetRouter(&pmsg.LeaveRequest{}, game.ChanRPC)       // 离开房间注册
	msg.Processor.SetRouter(&pmsg.MatchRoomRequest{}, game.ChanRPC)   // 房间匹配请求
	msg.Processor.SetRouter(&pmsg.MatchMemberRequest{}, game.ChanRPC) // 成员匹配请求

	msg.Processor.SetRouter(&pmsg.MemberReadyRequest{}, game.ChanRPC)       // 房间准备请求
	msg.Processor.SetRouter(&pmsg.MemberReadyCancelRequest{}, game.ChanRPC) // 成员取消准备请求

	msg.Processor.SetRouter(&pmsg.MatchRoomCancelRequest{}, game.ChanRPC) // 成员匹配请求
	msg.Processor.SetRouter(&pmsg.Answer{}, game.ChanRPC)                 // 回答问题请求注册
	msg.Processor.SetRouter(&pmsg.OverRequest{}, game.ChanRPC)            // 房间解散注册
	msg.Processor.SetRouter(&pmsg.Question{}, game.ChanRPC)               // 问题注册
	msg.Processor.SetRouter(&pmsg.MemberReliveRequest{}, game.ChanRPC)    // 复活注册
	msg.Processor.SetRouter(&pmsg.RoomInfoGetRequest{}, game.ChanRPC)     // 房间信息主动获取
	msg.Processor.SetRouter(&msg.TableCount{}, game.ChanRPC)              // 表格统计注册
	msg.Processor.SetRouter(&msg.DataUpdate{}, game.ChanRPC)              // 数据存储更新
	msg.Processor.SetRouter(&msg.DataRequest{}, game.ChanRPC)             // 数据存储更新
	msg.Processor.SetRouter(&msg.LoginRequest{}, game.ChanRPC)            // 数据存储更新

	msg.Processor.SetRouter(&pmsg.MemberLevelChange{}, game.ChanRPC) // 等级改变消息

}
