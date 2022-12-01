package handler

import (
	"catLog"
	"errors"
	"net"
	"net/http"
	"proto/msg"
	"remotemsg"
	"server"
	"server/client"
	"sort"
	"strings"

	"github.com/ip2location/ip2location-go"
	"google.golang.org/grpc"
)

type RankService struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.HelloClient
	S           *server.Server
}

func NewRank() *RankService {
	s := RankService{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

var RankInstance *RankService = NewRank()

var rankInfo = &RankInfo{}

var DB *ip2location.DB

func (s *RankService) Run(port string) {
	rankInfo.WorldRank = []*Rank{}
	rankInfo.CountryRank = make(map[string]*[]*Rank)
	rankInfo.ProvinceRank = make(map[string]*[]*Rank)
	rankInfo.CityRank = make(map[string]*[]*Rank)

	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
	if err != nil {
		catLog.Log(err)
		return
	}
	DB = db

	go func() {
		// 注册消息
		c1 := make(chan client.Msg)
		handler1 := &client.MsgHandler{}
		handler1.Chan = c1
		handler1.MsgID = remotemsg.RANKPULL
		client.RegisterHandler(handler1)

		c2 := make(chan client.Msg)
		handler2 := &client.MsgHandler{}
		handler2.Chan = c2
		handler2.MsgID = remotemsg.RANKUPDATE
		client.RegisterHandler(handler2)

		c3 := make(chan client.Msg)
		handler3 := &client.MsgHandler{}
		handler3.Chan = c3
		handler3.MsgID = remotemsg.RANKSELF
		client.RegisterHandler(handler3)

		for {
			select {
			case <-s.cat.Done:
				return
			case msg := <-c1:
				switch msg.Val.ID {
				case remotemsg.RANKPULL:
					pull(msg)
				case remotemsg.RANKUPDATE:
					update(msg)
				case remotemsg.RANKSELF:
					getSelf(msg)
				}
			case msg := <-c2:
				switch msg.Val.ID {
				case remotemsg.RANKPULL:
					pull(msg)
				case remotemsg.RANKUPDATE:
					update(msg)
				case remotemsg.RANKSELF:
					getSelf(msg)
				}
			case msg := <-c3:
				switch msg.Val.ID {
				case remotemsg.RANKPULL:
					pull(msg)
				case remotemsg.RANKUPDATE:
					update(msg)
				case remotemsg.RANKSELF:
					getSelf(msg)
				}
			}
		}
	}()

	// 注册消息
	// s.cat.Register(remotemsg.RANKPULL, pull)
	// s.cat.Register(remotemsg.RANKUPDATE, update)
	// s.cat.Register(remotemsg.RANKSELF, getSelf)
}

func getRankDataBy(ranks *[]*Rank, uid string) (*Rank, error) {
	for _, r := range *ranks {
		if r.UID == uid {
			return r, nil
		}
	}
	return nil, errors.New("找不到排行")
}

func handle(ranks *[]*Rank, count int, r *RequestRank) {

	// 如果数值比最后一名还小，则舍去更新
	if len(*ranks) >= count {
		lastR := (*ranks)[len((*ranks))-1]
		if lastR.Val > r.Val {
			return
		}
	}

	oldR, err := getRankDataBy(ranks, r.UID)

	if err != nil {
		oldR = &Rank{}
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

func update(data client.Msg) {
	rt := data.Client.R
	r := &RequestRank{}
	data.Val.ParseData(r)
	ip, _ := GetClientIp(rt)
	results, _ := DB.Get_all(ip)
	r.Country = results.Country_long
	r.CountryShort = results.Country_short
	// 更新全服排行数据
	handle(&rankInfo.WorldRank, 50, r)
	state := &upScore{}
	state.State = true
	catLog.Log("update排行", state)
	data.Client.MsgChan <- &client.BackMsg{
		MsgID: remotemsg.RANKPULL,
		Val:   state,
	}

}

func pull(data client.Msg) {
	backInfo := &BackRankInfo{}
	backInfo.WorldRank = rankInfo.WorldRank
	catLog.Log("拉取排行的消息", backInfo)
	data.Client.MsgChan <- &client.BackMsg{
		MsgID: remotemsg.RANKUPDATE,
		Val:   backInfo,
	}
}

func getSelf(data client.Msg) {
	rt := data.Client.R
	r := &RequestRank{}
	ip, _ := GetClientIp(rt)
	results, _ := DB.Get_all(ip)
	r.Country = results.Country_long
	r.CountryShort = results.Country_short

	oldR := &Rank{}
	oldR.NickName = r.NickName
	oldR.UID = r.UID
	oldR.Icon = r.Icon
	oldR.Val = r.Val
	oldR.Country = r.Country
	oldR.CountryShort = r.CountryShort
	catLog.Log("获取自己的消息", oldR)
	data.Client.MsgChan <- &client.BackMsg{
		MsgID: remotemsg.RANKSELF,
		Val:   oldR,
	}
}

// Rank
type Rank struct {
	UID          string `json:"uid"`
	Val          int    `json:"value"`
	Icon         string `json:"icon"`
	NickName     string `json:"nickname"`
	Country      string `json:"country"`
	CountryShort string `json:"countryShort"`
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

type upScore struct {
	State bool `json:"state"`
}

func GetClientIp(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
