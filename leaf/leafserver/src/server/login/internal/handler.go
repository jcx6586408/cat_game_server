package internal

import (
	"config"
	"io/ioutil"
	"leafserver/src/server/msg"
	"math/rand"
	"net/http"
	pmsg "proto/msg"
	"reflect"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

var wxConf *config.Config

func init() {
	handleMsg(&msg.WXCode{}, login)
	handleMsg(&pmsg.RoomPreAddRequest{}, preCreateRoom)

	wxConf = config.Read()
}

func login(args []interface{}) {
	wxcode := args[0].(*msg.WXCode)
	a := args[1].(gate.Agent)
	resp, _ := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + wxConf.Wx.Appid +
		"&secret=" + wxConf.Wx.AppSecret +
		"&js_code=" + wxcode.Code +
		"&grant_type=authorization_code")

	body, _ := ioutil.ReadAll(resp.Body)

	a.WriteMsg(&msg.WXcodeReply{
		Openid: string(body),
	})
}

func preCreateRoom(args []interface{}) {
	ran := rand.Intn(len(wxConf.Urls))
	a := args[1].(gate.Agent)
	url := wxConf.Urls[ran]
	log.Debug("下发服务器的路径: %v", url)
	a.WriteMsg(&pmsg.RoomPreAddReply{
		Url: url,
	})
}
