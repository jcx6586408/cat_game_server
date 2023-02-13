package internal

import (
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"
	pmsg "proto/msg"
	"reflect"
	"sort"

	"github.com/name5566/leaf/gate"
)

func init() {

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
	handler(&pmsg.MemberReadyRequest{}, roomMemberReady)
	handler(&pmsg.MemberReadyCancelRequest{}, roomMemberReadyCancel)
	handler(&msg.RoomStartPlay{}, roomStartPlay)
	handler(&pmsg.MemberReliveRequest{}, roomMatchMemberRelive)
	handler(&pmsg.RoomInfoGetRequest{}, roomInfoGet)
	handler(&pmsg.MemberLevelChange{}, roomMemberLevelChange)
	handler(&pmsg.Say{}, roomSay)
	handler(&msg.Ping{}, Hearbeat)
	handler(&msg.TableCount{}, tableCount)
	handler(&msg.DataUpdate{}, dataUpdate)
	handler(&msg.DataRequest{}, dataRequest)
	handler(&msg.LoginRequest{}, loginRequst)
	handler(&msg.QuestionLibRequest{}, getQuestionLib)
}

func getQuestionLib(args []interface{}) {
	req := args[0].(*msg.QuestionLibRequest)
	a := args[1].(gate.Agent)
	q := GetTestQuestion(req.Type, req.Level)
	a.WriteMsg(q)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func dataUpdate(args []interface{}) {
	req := args[0].(*msg.DataUpdate)
	// a := args[1].(gate.Agent)
	u, ok := Users[req.Uuid]
	if ok {
		u.UpdateMap(req.Key, req.Value)
	}
}

func dataRequest(args []interface{}) {
	req := args[0].(*msg.DataRequest)
	a := args[1].(gate.Agent)
	u, ok := Users[req.Uuid]
	if ok {
		val := u.GetData(req.Key)
		a.WriteMsg(&msg.DataReply{
			Key:   req.Key,
			Value: val,
		})
	}
}

func tableCount(args []interface{}) {
	req := args[0].(*msg.TableCount)
	a := args[1].(gate.Agent)
	back := &msg.TableGet{Questions: []*msg.QuestionCount{}}
	total := []*msg.QuestionCount{}
	for _, v := range Questions.QuestionMap {
		total = append(total, &msg.QuestionCount{ID: v.Q.ID, Win: v.win, Fail: v.fail})
	}
	sort.Slice(total, func(i, j int) bool {
		return total[i].ID > total[j].ID
	})
	if req.Min >= len(total) {
		a.WriteMsg(back)
		return
	}
	if req.Max >= len(total) {
		back.Questions = total[req.Min:]
		a.WriteMsg(back)
		return
	}
	back.Questions = total[req.Min:req.Max]
	a.WriteMsg(back)
}

func login(args []interface{}) {
	a := args[1].(gate.Agent)
	wxcode := &msg.WXCode{}
	resp, _ := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + ServerConf.Wx.Appid +
		"&secret=" + ServerConf.Wx.AppSecret +
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
