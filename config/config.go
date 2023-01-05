package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Rank      *Rank    `json:"rank"`      // 排行榜配置
	Wx        *Wx      `json:"wx"`        // 微信openid配置
	DB        *DB      `json:"DB"`        // 数据库连接路径
	Urls      []string `json:"Urls"`      // 游戏房间服务器路径
	SeverType string   `json:"SeverType"` // 游戏服类型（网关服gate，中心服center，房间服room）
}

type DB struct {
	User     string
	Password string
	Host     string
	DB       string
	MinAlive int
	MaxAlive int
	MaxIdle  int
}

type Wx struct {
	Appid     string `json:"appid"`
	AppSecret string `json:"appSecret"`
}

type Rank struct {
	Port              string `json:"Port"`
	RedisUrl          string `json:"RedisUrl"`
	RedisPassword     string `json:"RedisPassword"`
	WorldRankCount    int    `json:"worldRankCount"`
	CountryRankCount  int    `json:"countryRankCount"`
	ProvinceRankCount int    `json:"provinceRankCount"`
	CityRankCount     int    `json:"cityRankCount"`
	Max               int    `json:"Max"`
}

func Read() *Config {
	conf := &Config{}
	data, err := ioutil.ReadFile("./conf/server.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		panic("解析json文件出错")
	}
	return conf
}
