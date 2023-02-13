package rank

import (
	"config"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"
	"sort"
	"sync"

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

var (
	DB *ip2location.DB

	rankInfo = &RankInfo{}

	Conf *config.Config

	IPLocationPath string

	serversMax int

	curServer int

	lock sync.RWMutex
)

func RankInit() {
	Conf = config.Read()
	serversMax = len(Conf.Urls)
	curServer = 0
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

func handleDele(l *sync.RWMutex, ranks *[]*msg.Rank, count int, r *msg.Rank) {
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
		oldR.Skin = r.Skin
	}

	// 更新排行
	if len(*ranks) > 1 {
		sort.SliceStable((*ranks), func(i, j int) bool {
			return (*ranks)[i].Val > (*ranks)[j].Val
		})

		(*ranks) = delete((*ranks), oldR)

		if len((*ranks)) > count {
			(*ranks) = (*ranks)[:count]
		}
	}
}

func handle(l *sync.RWMutex, ranks *[]*msg.Rank, count int, r *msg.Rank) {
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
		oldR.Skin = r.Skin
		(*ranks) = append((*ranks), oldR) // 加入排行
	} else {
		oldR.Val = r.Val
		oldR.NickName = r.NickName
		oldR.Icon = r.Icon
		oldR.City = r.City
		oldR.Skin = r.Skin
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

	lock.RLock()
	// 城市
	cityRanks, ok := rankInfo.CityRank[r.City]
	lock.RUnlock()
	if !ok {
		cityRanks = &([]*msg.Rank{})
	}
	lock.Lock()
	handle(&lock, cityRanks, Conf.Rank.CityRankCount, r) // 城市排行更新
	rankInfo.CityRank[r.City] = cityRanks                // 重新赋值
	lock.Unlock()
	return c.JSON(http.StatusOK, r)
}

func RankDele(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	// 城市
	lock.RLock()
	cityRanks, ok := rankInfo.CityRank[r.City]
	lock.RUnlock()
	if !ok {
		cityRanks = &([]*msg.Rank{})
	}
	lock.Lock()
	handleDele(&lock, cityRanks, Conf.Rank.CityRankCount, r) // 城市排行更新
	rankInfo.CityRank[r.City] = cityRanks                    // 重新赋值
	lock.Unlock()
	return c.JSON(http.StatusOK, r)
}

func RankPull(c echo.Context) error {
	return c.JSON(http.StatusOK, rankInfo.WorldRank)
}

func RankCityPull(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	lock.RLock()
	backInfo, ok := rankInfo.CityRank[r.City]
	lock.RUnlock()
	if !ok {
		return c.JSON(http.StatusOK, "")
	}

	var max int = len(*backInfo)
	if len(*backInfo) > 35 {
		max = 35
	}
	return c.JSON(http.StatusOK, (*backInfo)[0:max])
}

func GetSelf(c echo.Context) error {
	oldR := &msg.Rank{}
	return c.JSON(http.StatusOK, oldR)
}

func RoomCreate(c echo.Context) error {
	url := Conf.Urls[curServer]
	curServer++
	if curServer >= serversMax {
		curServer = 0
	}
	return c.JSON(http.StatusOK, &pmsg.RoomPreAddReply{
		Url: url,
	})
}

type WXCode struct {
	Code string
}

type BytedanceCode struct {
}

type OpenID struct {
	Openid string `json:"openid"`
}

// 微信
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
	openid := &OpenID{}
	json.Unmarshal(body, openid)
	return c.JSON(http.StatusOK, openid)
}

// 字节跳动
func GetBytedanceOpenID(c echo.Context) error {
	wxcode := &WXCode{}
	ParseNetBody(wxcode, c.Request().Body)
	resp, err := http.Get("https://minigame.zijieapi.com/mgplatform/api/apps/jscode2session?appid=" + Conf.Bytedance.Appid +
		"&secret=" + Conf.Bytedance.AppSecret +
		"&code=" + wxcode.Code +
		"&anonymous_code=")
	if err != nil {
		return c.String(http.StatusOK, "")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	openid := &OpenID{}
	json.Unmarshal(body, openid)
	return c.JSON(http.StatusOK, openid)
}

func ParseNetBody(i interface{}, r io.ReadCloser) {
	d, _ := ioutil.ReadAll(r)
	json.Unmarshal(d, i)
}

func delete(a []*msg.Rank, elem *msg.Rank) []*msg.Rank {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
