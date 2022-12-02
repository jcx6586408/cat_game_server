package rank

import (
	"github.com/ip2location/ip2location-go"
	"github.com/name5566/leaf/gate"
)

// Rank
type Rank struct {
	UID          string `json:"uid"`
	Val          int    `json:"value"`
	Icon         string `json:"icon"`
	NickName     string `json:"nickname"`
	Country      string `json:"country"`
	CountryShort string `json:"countryShort"`
}

type RankSelfRequest struct {
	UID      string `json:"uid"`
	Val      int    `json:"value"`
	Icon     string `json:"icon"`
	NickName string `json:"nickname"`
}

type RequestRank struct {
	Rank
}

type BackRankInfo struct {
	WorldRank []*Rank
}

type RankInfo struct {
	CountryRank  map[string]*[]*Rank
	WorldRank    []*Rank
	ProvinceRank map[string]*[]*Rank
	CityRank     map[string]*[]*Rank
}

type UpScore struct {
	State bool `json:"state"`
}

var DB *ip2location.DB

var rankInfo = &RankInfo{}

func init() {
	rankInfo.WorldRank = []*Rank{}
	rankInfo.CountryRank = make(map[string]*[]*Rank)
	rankInfo.ProvinceRank = make(map[string]*[]*Rank)
	rankInfo.CityRank = make(map[string]*[]*Rank)
}

func GetSelf(args []interface{}) {
	// 收到的 Hello 消息
	r := args[0].(*RankSelfRequest)
	// 消息的发送者
	a := args[1].(gate.Agent)
	results, _ := DB.Get_all(a.RemoteAddr().String())
	oldR := &Rank{}
	oldR.NickName = r.NickName
	oldR.UID = r.UID
	oldR.Icon = r.Icon
	oldR.Val = r.Val
	oldR.Country = results.Country_long
	oldR.CountryShort = results.Country_short
	a.WriteMsg(oldR)
}
