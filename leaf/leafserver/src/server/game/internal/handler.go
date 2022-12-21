package internal

import (
	"config"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"
	pmsg "proto/msg"
	"reflect"
	"storage/redis"

	"github.com/name5566/leaf/gate"
)

func init() {
	// 数据库连接
	// storage.Connect()
	redis.ConnectReids()
	// 模块初始化
	RankInit()
	ConstInit()
	// RoomManagerInit()
	wxConf = config.Read()
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.RankSelfRequest{}, GetSelf)
	handler(&msg.Rank{}, RankUpdate)
	handler(&msg.RankPull{}, RankPull)
	handler(&msg.WXCode{}, login)

	// 房间消息注册
	handler(&pmsg.CreateRoomRequest{}, roomCreate)
	handler(&pmsg.AddRequest{}, roomAdd)
	handler(&pmsg.LeaveRequest{}, roomLeave)
	handler(&pmsg.Answer{}, roomAnswer)
	handler(&pmsg.OverRequest{}, roomOver)
	handler(&pmsg.MatchRoomRequest{}, roomMatch)
	handler(&pmsg.MatchMemberRequest{}, roomMatchMember)
	handler(&pmsg.MatchRoomCancelRequest{}, roomMatchCanel)
	handler(&pmsg.MatchMemberCancelRequest{}, roomMatchMemberCanel)
	handler(&msg.RoomStartPlay{}, roomStartPlay)
	handler(&pmsg.MemberReliveRequest{}, roomMatchMemberRelive)
	handler(&pmsg.RoomInfoGetRequest{}, roomInfoGet)
	handler(&msg.TableCount{}, tableCount)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func tableCount(args []interface{}) {
	req := args[0].(*msg.TableCount)
	a := args[1].(gate.Agent)
	a.WriteMsg(redis.GetWinTableRank(req.Min, req.Max))
	a.WriteMsg(redis.GetFailTableRank(req.Min, req.Max))
}

var wxConf *config.Config

func login(args []interface{}) {
	a := args[1].(gate.Agent)
	wxcode := &msg.WXCode{}
	resp, _ := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + wxConf.Wx.Appid +
		"&secret=" + wxConf.Wx.AppSecret +
		"&js_code=" + wxcode.Code +
		"&grant_type=authorization_code")

	body, _ := ioutil.ReadAll(resp.Body)
	// guid := AgentUsers[a]
	openid := string(body)
	// loginHandle(guid, openid) // 登录处理
	a.WriteMsg(&msg.WXcodeReply{
		Openid: openid,
	})
}
