package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Gate    *Gate    `json:"gate"`    // 中心服配置
	Excel   *Excel   `json:"excel"`   // 云表格配置
	Rank    *Rank    `json:"rank"`    // 排行榜配置
	Crt     *Crt     `json:"crt"`     // 证书配置
	Wx      *Wx      `json:"wx"`      // 微信openid配置
	Storage *Storage `json:"storage"` // 游戏存储
	Game    *Game    `json:"game"`    // 游戏逻辑服
	DB      string   `json:"DB"`      // 数据库连接路径
}

type Base struct {
	ID   int    `json:"ID"` // 模块ID
	Port string `json:"port"`
}

type Gate struct {
	Base
	HttpsPort    string   `json:"httpsPort"`
	ListenerType []string `json:"listenerType"`
}

type Game struct {
	Base
}

type Storage struct {
	Base
}

type Wx struct {
	Base
	Appid     string `json:"appid"`
	AppSecret string `json:"appSecret"`
}

type Crt struct {
	Base
	CertFile string `json:"crt"`
	KeyFile  string `json:"crtKey"`
}

type Excel struct {
	Base
}

type Rank struct {
	Base
	WorldRankCount    int `json:"worldRankCount"`
	CountryRankCount  int `json:"countryRankCount"`
	ProvinceRankCount int `json:"provinceRankCount"`
	CityRankCount     int `json:"cityRankCount"`
}

func (c *Config) GetModules() *(map[int]*Base) {
	bs := make(map[int]*Base)
	bs[c.Excel.ID] = &c.Excel.Base
	bs[c.Rank.ID] = &c.Rank.Base
	bs[c.Crt.ID] = &c.Crt.Base
	bs[c.Game.ID] = &c.Game.Base
	bs[c.Storage.ID] = &c.Storage.Base
	bs[c.Wx.ID] = &c.Wx.Base
	return &bs
}

func Read() *Config {
	conf := &Config{}
	data, err := ioutil.ReadFile("./server.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		panic("解析json文件出错")
	}
	return conf
}
