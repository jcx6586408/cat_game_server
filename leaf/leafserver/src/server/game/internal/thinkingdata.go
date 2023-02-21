package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

var (
	appid = "5a7a9d5530304eb788aa62ae005510b5"
)

func TgaTrack(uid int64, action string, params map[string]interface{}) {
	now := time.Now()

	msg := make(map[string]interface{})
	msg["#account_id"] = uid
	msg["#type"] = "track"
	msg["#event_name"] = action
	msg["#time"] = now.Format("2006-01-02 15:04:05")

	if nil != params {
	} else {
		params = make(map[string]interface{})
	}
	params["uid"] = uid
	params["sendTime"] = msg["#time"]
	now1 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)
	params["send_min_t"] = now1.Format("2006-01-02 15:04:05")
	msg["properties"] = params

	// beego.Debug(msg["#time"])

	body, _ := json.Marshal(&msg)

	res, err := http.PostForm("https://global-receiver-ta.thinkingdata.cn",
		url.Values{"appid": {"5a7a9d5530304eb788aa62ae005510b5"}, "data": {string(body)}})
	if nil != err {
		// beego.Error("发送tga事件失败, err:", err.Error())
		return
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	// beego.Debug(fmt.Sprintf("发送tga事件结果: %s", buf.String()))
}

// 设置用户属性
func TgaUserSet(uid int64, params map[string]interface{}) {

	if nil == params {
		return
	}

	now := time.Now()

	msg := make(map[string]interface{})
	msg["#account_id"] = uid
	msg["#type"] = "user_set"
	msg["#time"] = now.Format("2006-01-02 15:04:05")
	msg["properties"] = params

	body, _ := json.Marshal(&msg)
	res, err := http.PostForm("https://tga.hortorgames.com/sync_data",
		url.Values{"appid": {appid}, "data": {string(body)}})
	if nil != err {
		// beego.Error("设置tga用户属性失败, err:", err.Error())
		return
	}
	defer res.Body.Close()
}

// 设置用户属性
func TgaUserSetOnce(uid int64, params map[string]interface{}) {
	if nil == params {
		return
	}

	now := time.Now()

	msg := make(map[string]interface{})
	msg["#account_id"] = uid
	msg["#type"] = "user_setOnce"
	msg["#time"] = now.Format("2006-01-02 15:04:05")
	msg["properties"] = params

	body, _ := json.Marshal(&msg)
	res, err := http.PostForm("https://tga.hortorgames.com/sync_data",
		url.Values{"appid": {appid}, "data": {string(body)}})
	if nil != err {
		return
	}
	defer res.Body.Close()
}

func TrackCheatAct(uid int64, kind, step string, lvl int) {

	params := make(map[string]interface{})
	if "" != kind {
		params["kind"] = kind
	}
	if "" != step {
		params["step"] = step
	}
	if 0 != lvl {
		params["lvl"] = lvl
	}
	go TgaTrack(uid, "cheat", params)
}
