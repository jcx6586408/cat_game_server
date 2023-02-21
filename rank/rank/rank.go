package rank

import (
	"config"
	"context"
	"encoding/json"
	"errors"
	"excel"
	"fmt"
	"io"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"
	"sort"
	"strconv"
	"sync"

	pmsg "proto/msg"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
)

type Ranks map[string]*[]*msg.Rank

type RequestRank struct {
	msg.Rank
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

type BackRanks struct {
	Ranks []*msg.Rank
	Total int
}

var (
	rankInfo = &RankInfo{}

	Conf *config.Config

	serversMax int

	curServer int

	lock sync.RWMutex

	// Rdb redis.Pipeliner
	Rdb *redis.Client

	Pipe redis.Pipeliner
	Max  int64

	ctx = context.Background()

	tables *excel.ExcelManager

	LevelLib []*Level

	LevelDBLib       map[string]*LevelDB
	LevelDbLibBylv   map[string]int
	LevelDbLibByName map[int]string

	countName = "Count"
)

func RankInit() {
	LevelDBLib = make(map[string]*LevelDB)
	LevelDbLibBylv = make(map[string]int)
	LevelDbLibByName = make(map[int]string)
	Conf = config.Read()
	tables = excel.Read()
	LevelLib = ToLevelLib()
	serversMax = len(Conf.Urls)
	curServer = 0
	rankInfo.WorldRank = []*msg.Rank{}
	rankInfo.CountryRank = make(map[string]*[]*msg.Rank)
	rankInfo.ProvinceRank = make(map[string]*[]*msg.Rank)
	rankInfo.CityRank = make(map[string]*[]*msg.Rank)

	Max = int64(Conf.Rank.Max)
	println("最大排行数量: ", Max)
	fmt.Println("redis连接参数 :", Conf.Rank.RedisUrl, Conf.Rank.RedisPassword)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Conf.Rank.RedisUrl,
		Password: Conf.Rank.RedisPassword,
		DB:       0,
		PoolSize: 1000,
	})
	result, err := Rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis数据库连接错误 :", err)
		return
	}
	println(result)
	fmt.Printf("******************redis数据库连接成功******************\n")
	sort.SliceStable((LevelLib), func(i, j int) bool {
		return (LevelLib)[i].ID < (LevelLib)[j].ID
	})

	// 读取redis数据
	for _, v := range LevelLib {
		val, err := LoadFromDb(v.Name)
		countVal, countErr := LoadFromDb(v.Name + countName)
		LevelDbLibBylv[v.Name] = v.ID
		LevelDbLibByName[v.ID] = v.Name
		println("段位信息: ", v.ID, v.Name)
		if err == nil {
			var city = []*msg.Rank{}
			err := json.Unmarshal([]byte(val), &city)
			if err == nil {
				rankInfo.CityRank[v.Name] = &city
			}
		}
		if countErr == nil {
			vv, err := strconv.Atoi(countVal)
			LevelDBLib[v.Name].Count = vv
			if err != nil {
				println("解析统计人数错误")
			}
		}
	}
}

func SaveCountToDB(key string, add int) {
	var uidKey = key + countName
	r, e := LoadFromDb(uidKey)
	if e == nil {
		rr, ee := strconv.Atoi(r)
		if ee == nil {
			rr = rr + add
			Rdb.Set(ctx, uidKey, rr, 0)
		}
	} else {
		if e == redis.Nil {
			Rdb.Set(ctx, uidKey, 1, 0)
		}
	}
}

func GetCountFromDB(key string) int {
	var uidKey = key + countName
	r, e := LoadFromDb(uidKey)
	if e != nil {
		return 0
	}
	rr, ee := strconv.Atoi(r)
	if ee != nil {
		return 0
	}
	return rr
}

func SaveToDB(key string) {
	arr, ok := rankInfo.CityRank[key]
	if !ok {
		return
	}
	value, e := json.Marshal(&arr)
	if e != nil {
		return
	}
	var st = Rdb.Set(ctx, key, string(value), 0)
	err := st.Err()
	if err != nil {
		fmt.Println("set err :", err)
		return
	}
}

func LoadFromDb(key string) (string, error) {
	var result = Rdb.Get(ctx, key)
	return result.Result()
}

func getRankDataBy(ranks *[]*msg.Rank, uid string) (*msg.Rank, error) {
	for _, r := range *ranks {
		if r.UID == uid {
			return r, nil
		}
	}
	return nil, errors.New("找不到排行")
}

func GetNextCity(city string) string {
	return LevelDbLibByName[LevelDbLibBylv[city]+1]
}

func handleDele(l *sync.RWMutex, ranks *[]*msg.Rank, count int, r *msg.Rank) {
	// 如果数值比最后一名还小，则舍去更新
	if len(*ranks) >= count {
		lastR := (*ranks)[len((*ranks))-1]
		if lastR.Val >= r.Val {
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
	} else {

		(*ranks) = delete((*ranks), oldR)

	}
	SaveToDB(r.City)
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
	SaveToDB(r.City)
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
	// 更新redis
	if r.City == "" {
		var cur = LevelDbLibByName[1]
		fmt.Printf("\n%v\n%v\n", cur, LevelDbLibByName)
		SaveCountToDB(cur, 1)
		return c.JSON(http.StatusOK, r)
	}
	var next = GetNextCity(r.City)
	SaveCountToDB(r.City, -1)
	SaveCountToDB(next, 1)
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

func RankCityPull_V1(c echo.Context) error {
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

	return c.JSON(http.StatusOK, &BackRanks{
		Ranks: (*backInfo)[0:max],
		Total: GetCountFromDB(r.City),
	})
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
