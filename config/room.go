package config

import (
	"encoding/json"
	"io/ioutil"
)

type RoomConfig struct {
	PrepareTime   int // 比赛准备时间
	AnswerTime    int // 单次回答问题时间
	MaxMember     int // 房间成员数量
	QuestionCount int // 答题数量
}

func ReadRoom() *RoomConfig {
	conf := &RoomConfig{}
	data, err := ioutil.ReadFile("./room.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		panic("解析json文件出错")
	}
	return conf
}
