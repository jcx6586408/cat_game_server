package redis

import (
	"config"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/ip2location/ip2location-go"
	"github.com/labstack/echo"
)

var (
	Rdb *redis.Client

	worldRank = "world"

	win  = "win"
	fail = "fail"

	DB *ip2location.DB

	conf *config.Config
	Pipe redis.Pipeliner
	Max  int64

	ctx = context.Background()
)

func ConnectReids() {
	conf = config.Read()
	Max = int64(conf.Rank.Max)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Rank.RedisUrl,
		Password: "",
		DB:       0,
	})

	// Pipe = Rdb.Pipeline()
	println("******************redis数据库连接成功******************")
	iplocationInit()
}

func GetAddr() string {
	return conf.Rank.Port
}

func iplocationInit() {
	db, err := ip2location.OpenDB("../../../.././IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
	if err != nil {
		// catLog.Log(err)
		db, err = ip2location.OpenDB("./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN")
		if err != nil {
			return
		}
	}
	DB = db
}

func AddWinTable(uid string, score float64) {
	AddTable(win, uid, score)
}

func AddFailTable(uid string, score float64) {
	AddTable(fail, uid, score)
}

// 加入统计
func AddTable(key, uid string, score float64) {
	Rdb.ZIncr(ctx, key, &redis.Z{Member: uid, Score: score})
}

func GetWinTableRank(min, max int) *msg.BackTable {
	arr := Rdb.ZRevRangeWithScores(ctx, win, int64(min), int64(max)).Val()
	return &msg.BackTable{
		Arr:  arr,
		Name: win,
	}
}

func GetFailTableRank(min, max int) *msg.BackTable {
	arr := Rdb.ZRevRangeWithScores(ctx, fail, int64(min), int64(max)).Val()
	return &msg.BackTable{
		Arr:  arr,
		Name: fail,
	}
}

// 加入排行榜
func AddRank(key, uid string, score float64) {
	Rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: uid})
	count := Rdb.ZCard(ctx, key)
	if count.Val() > Max {
		Rdb.ZRemRangeByRank(ctx, key, Max, -1)
	}
}

// 获取自身世界排行
func GetSelfWorldRank(uid string) int64 {
	return Rdb.ZRevRank(ctx, worldRank, uid).Val() + 1
}

// 获取自身城市排行
func GetSelfCityRank(city, uid string) int64 {
	return Rdb.ZRevRank(ctx, city, uid).Val() + 1
}

// 更新排行榜
func UpdateRank(city, uid string, score float64) {
	// 世界排行更新
	f := Rdb.ZScore(ctx, worldRank, uid).Val()
	if f < score {
		AddRank(worldRank, uid, score)
	}

	// 城市排行更新
	f = Rdb.ZScore(ctx, city, uid).Val()
	if f < score {
		AddRank(city, uid, score)
	}
}

// 获取排行榜
func GetRank(key string, start int64, end int64) []redis.Z {
	return Rdb.ZRevRangeWithScores(ctx, key, start, end).Val()
}

// 获取世界排行
func GetWorldRank(start int64, end int64) []redis.Z {
	return GetRank(worldRank, start, end)
}

// 获取城市排行
func GetCityRank(city string, start, end int64) []redis.Z {
	return GetRank(city, start, end)
}

func RankUpdate(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	r.City = results.City
	UpdateRank(r.City, r.UID, float64(r.Val))
	return c.JSON(http.StatusOK, r)
}

func RankPull(c echo.Context) error {
	return c.JSON(http.StatusOK, GetWorldRank(0, int64(conf.Rank.WorldRankCount)))
}

func RankCityPull(c echo.Context) error {
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	backInfo := GetCityRank(results.City, 0, int64(conf.Rank.CityRankCount))
	return c.JSON(http.StatusOK, backInfo)
}

func GetSelf(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	addr := c.RealIP()
	results, _ := DB.Get_all(addr)
	oldR := &msg.Rank{}
	oldR.WorldRank = GetSelfWorldRank(r.UID)
	oldR.CityRank = GetSelfCityRank(results.City, r.UID)
	oldR.City = results.City
	return c.JSON(http.StatusOK, oldR)
}

func ParseNetBody(i interface{}, r io.ReadCloser) {
	d, _ := ioutil.ReadAll(r)
	json.Unmarshal(d, i)
}
