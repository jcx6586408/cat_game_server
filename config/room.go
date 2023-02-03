package config

import (
	"encoding/json"
	"io/ioutil"
)

var RoomConfPath string

type RoomConfig struct {
	PrepareTime           int // 比赛准备时间
	AnswerTime            int // 单次回答问题时间
	MaxMember             int // 房间成员数量
	MaxInvite             int // 邀请成员上限
	ReliveWaitTime        int // 房间复活等待时间
	QuestionCountMinLimit int
	QuestionCountDur      int
	RobotActionMin        float32
	RobotActionMax        float32
}

func ReadRoom() *RoomConfig {
	conf := &RoomConfig{}
	data, err := ioutil.ReadFile(RoomConfPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		panic("解析json文件出错")
	}
	return conf
}
