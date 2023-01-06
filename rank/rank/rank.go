package rank

import (
	"config"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"leafserver/src/server/msg"
	"math/rand"
	"net/http"
	"sort"

	pmsg "proto/msg"

	"github.com/ip2location/ip2location-go"
	"github.com/labstack/echo"
)

type Ranks map[string]*[]*msg.Rank

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

var Conf *config.Config

var IPLocationPath string

func RankInit() {
	Conf = config.Read()

	db, err := ip2location.OpenDB("../../../.././IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
	if err != nil {
		// catLog.Log(err)
		db, err = ip2location.OpenDB(IPLocationPath)
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
		oldR.City = r.City
		(*ranks) = append((*ranks), oldR) // 加入排行
	} else {
		oldR.Val = r.Val
		oldR.NickName = r.NickName
		oldR.Icon = r.Icon
		oldR.City = r.City
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

func RankUpdate(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	r.Country = results.Country_long
	r.CountryShort = results.Country_short
	r.City = results.City
	handle(&rankInfo.WorldRank, Conf.Rank.WorldRankCount, r) // 全服排行更新

	// 城市
	cityRanks, ok := rankInfo.CityRank[r.City]
	if !ok {
		cityRanks = &([]*msg.Rank{})
	}
	handle(cityRanks, Conf.Rank.CityRankCount, r) // 城市排行更新
	rankInfo.CityRank[r.City] = cityRanks         // 重新赋值
	return c.JSON(http.StatusOK, r)
}

func RankPull(c echo.Context) error {
	return c.JSON(http.StatusOK, rankInfo.WorldRank)
}

func RankCityPull(c echo.Context) error {
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	backInfo := rankInfo.CityRank[results.City]
	return c.JSON(http.StatusOK, backInfo)
}

func GetSelf(c echo.Context) error {
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	oldR := &msg.Rank{}
	oldR.Country = results.Country_long
	oldR.CountryShort = results.Country_short
	return c.JSON(http.StatusOK, oldR)
}

func RoomCreate(c echo.Context) error {
	ran := rand.Intn(len(Conf.Urls))
	url := Conf.Urls[ran]
	return c.JSON(http.StatusOK, &pmsg.RoomPreAddReply{
		Url: url,
	})
}

type WXCode struct {
	Code string
}

func GetOpenID(c echo.Context) error {
	wxcode := &WXCode{}
	ParseNetBody(wxcode, c.Request().Body)
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + Conf.Wx.Appid +
		"&secret=" + Conf.Wx.AppSecret +
		"&js_code=" + wxcode.Code +
		"&grant_type=authorization_code")
	if err != nil {
		return c.String(http.StatusOK, "")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return c.String(http.StatusOK, string(body))
}

func ParseNetBody(i interface{}, r io.ReadCloser) {
	d, _ := ioutil.ReadAll(r)
	json.Unmarshal(d, i)
}
