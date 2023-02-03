package msg

import (
	pmsg "proto/msg"

	predis "github.com/go-redis/redis/v8"

	"github.com/name5566/leaf/network/json"
)

// "github.com/name5566/leaf/network"

// var Processor network.Processor
var Processor = json.NewProcessor()

// var ProbufProcessor = protobuf.NewProcessor()

func init() {
	// 登录消息
	Processor.Register(&Login{})
	Processor.Register(&LoginRequest{})
	// ProbufProcessor.Register(&pmsg.RoomPreAddRequest{})
	// ProbufProcessor.Register(&pmsg.RoomPreAddReply{})
	// ProbufProcessor.Register(&pmsg.RoomInfoGetRequest{})

	// 排行榜消息
	Processor.Register(&RankSelfRequest{})
	Processor.Register(&Rank{})
	Processor.Register(&BackRankInfo{})
	Processor.Register(&RankPull{})

	// 房间消息注册
	Processor.Register(&pmsg.RoomPreAddRequest{})  // 房间预创建
	Processor.Register(&pmsg.RoomPreAddReply{})    // 房间预创建返回
	Processor.Register(&pmsg.RoomInfoGetRequest{}) // 房间信息主动请求

	Processor.Register(&pmsg.CreateRoomRequest{}) // 房间创建
	Processor.Register(&pmsg.CreateRoomReply{})   // 房间创建返回

	Processor.Register(&pmsg.RoomTime{}) // 房间创建返回

	// 复活
	Processor.Register(&pmsg.MemberReliveRequest{}) // 成员复活
	Processor.Register(&pmsg.MemberReliveReply{})   // 成员复活返回

	Processor.Register(&RoomStartPlay{}) // 房间主动开始游戏

	Processor.Register(&pmsg.AddRequest{}) // 加入房间注册
	Processor.Register(&pmsg.AddReply{})   // 加入房间注册回复

	Processor.Register(&pmsg.LeaveRequest{}) // 离开房间注册
	Processor.Register(&pmsg.LeaveReply{})   // 离开房间注册回复

	Processor.Register(&pmsg.MatchRoomRequest{}) // 房间匹配请求
	Processor.Register(&pmsg.MatchRoomReply{})   // 房间匹配请求回复

	Processor.Register(&pmsg.MatchMemberRequest{}) // 成员匹配请求
	Processor.Register(&pmsg.MatchMemberReply{})   // 成员匹配请求回复

	Processor.Register(&pmsg.MemberReadyRequest{})       // 成员准备请求
	Processor.Register(&pmsg.MemberReadyCancelRequest{}) // 成员准备取消请求回复

	Processor.Register(&pmsg.MatchRoomCancelRequest{}) // 房间匹配取请求
	Processor.Register(&pmsg.MatchRoomCancelReply{})   // 房间匹配取消请求回复

	Processor.Register(&pmsg.MatchMemberCancelRequest{}) // 成员匹配取消请求
	Processor.Register(&pmsg.MatchMemberCancelReply{})   // 成员匹配取消请求回复

	Processor.Register(&pmsg.Answer{})         // 回答问题请求注册
	Processor.Register(&pmsg.AnswerEndReply{}) // 答题结束

	Processor.Register(&pmsg.OverRequest{})   // 房间解散注册
	Processor.Register(&pmsg.Question{})      // 问题注册
	Processor.Register(&pmsg.RoomInfoReply{}) // 比赛房间状态信息注册
	Processor.Register(&WXCode{})             // 微信请求
	Processor.Register(&WXcodeReply{})        // 微信请求
	Processor.Register(&pmsg.RoomAddFail{})   // 房间加入失败返回

	Processor.Register(&pmsg.MemberLevelChange{}) // 房间加入失败返回
	Processor.Register(&pmsg.Say{})               // 房间发送说话内容

	// 表格消息注册
	Processor.Register(&QuestionLibRequest{}) // 房间发送说话内容
	Processor.Register(&TableCount{})
	Processor.Register(&BackTable{})
	Processor.Register(&DataUpdate{})
	Processor.Register(&DataRequest{})
	Processor.Register(&DataReply{})
	Processor.Register(&TableGet{})
}

type WXCode struct {
	Code string
}

type WXcodeReply struct {
	Openid string
}

type TableCount struct {
	Min int
	Max int
}

type TableGet struct {
	Questions []*QuestionCount
}

type QuestionCount struct {
	ID   int
	Win  int
	Fail int
}

type BackTable struct {
	Arr  []predis.Z
	Name string
}

type DataUpdate struct {
	Uuid  string
	Key   string
	Value string
}

type DataRequest struct {
	Uuid string
	Key  string
}

type DataReply struct {
	Key   string
	Value string
}

type QuestionLibRequest struct {
	Type  string
	Level int
}
