package internal

import (
	"errors"
	"leafserver/src/server/msg"
	"sort"
	"strings"

	"github.com/ip2location/ip2location-go"
	"github.com/name5566/leaf/gate"
)

type RequestRank struct {
	msg.Rank
}

type BackRankInfo struct {
	WorldRank []*msg.Rank
}

type RankInfo struct {
	CountryRank  map[string]*[]*msg.Rank
	WorldRank    []*msg.Rank
	ProvinceRank map[string]*[]*msg.Rank
	CityRank     map[string]*[]*msg.Rank
}

type UpScore struct {
	State bool `json:"state"`
}

var DB *ip2location.DB

var rankInfo = &RankInfo{}

func RankInit() {
	db, err := ip2location.OpenDB("../../../.././IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
	if err != nil {
		// catLog.Log(err)
		db, err = ip2location.OpenDB("./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
		if err != nil {
			return
		}
	}
	DB = db
	rankInfo.WorldRank = []*msg.Rank{}
	rankInfo.CountryRank = make(map[string]*[]*msg.Rank)
	rankInfo.ProvinceRank = make(map[string]*[]*msg.Rank)
	rankInfo.CityRank = make(map[string]*[]*msg.Rank)
}

func GetSelf(args []interface{}) {
	// 收到的 Hello 消息
	r := args[0].(*msg.RankSelfRequest)
	// 消息的发送者
	a := args[1].(gate.Agent)
	addr := strings.Split(a.RemoteAddr().String(), ":")[0]
	results, _ := DB.Get_all(addr)
	oldR := &msg.Rank{}
	oldR.NickName = r.NickName
	oldR.UID = r.UID
	oldR.Icon = r.Icon
	oldR.Val = r.Val
	oldR.Country = results.Country_long
	oldR.CountryShort = results.Country_short
	a.WriteMsg(oldR)
}

func getRankDataBy(ranks *[]*msg.Rank, uid string) (*msg.Rank, error) {
	for _, r := range *ranks {
		if r.UID == uid {
			return r, nil
		}
	}
	return nil, errors.New("找不到排行")
}

func handle(ranks *[]*msg.Rank, count int, r *msg.Rank) {

	// 如果数值比最后一名还小，则舍去更新
	if len(*ranks) >= count {
		lastR := (*ranks)[len((*ranks))-1]
		if lastR.Val > r.Val {
			return
		}
	}

	oldR, err := getRankDataBy(ranks, r.UID)

	if err != nil {
		oldR = &msg.Rank{}
		oldR.NickName = r.NickName
		oldR.UID = r.UID
		oldR.Icon = r.Icon
		oldR.Val = r.Val
		oldR.Country = r.Country
		oldR.CountryShort = r.CountryShort
		(*ranks) = append((*ranks), oldR) // 加入排行
	} else {
		oldR.Val = r.Val
		oldR.NickName = r.NickName
		oldR.Icon = r.Icon
	}

	// 更新排行
	if len(*ranks) > 1 {
		sort.SliceStable((*ranks), func(i, j int) bool {
			return (*ranks)[i].Val > (*ranks)[j].Val
		})

		if len((*ranks)) > count {
			(*ranks) = (*ranks)[:count]
		}
	}
}

func RankUpdate(args []interface{}) {
	r := args[0].(*msg.Rank)
	// 消息的发送者
	a := args[1].(gate.Agent)
	addr := strings.Split(a.RemoteAddr().String(), ":")[0]
	results, _ := DB.Get_all(addr)
	r.Country = results.Country_long
	r.CountryShort = results.Country_short
	handle(&rankInfo.WorldRank, 50, r)
}

func RankPull(args []interface{}) {
	a := args[1].(gate.Agent)
	backInfo := &BackRankInfo{}
	backInfo.WorldRank = rankInfo.WorldRank
	a.WriteMsg(backInfo)
}
