package redis

import (
	"config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/ip2location/ip2location-go"
	"github.com/labstack/echo"
)

var (
	// Rdb redis.Pipeliner
	Rdb *redis.Client

	worldRank = "world"

	win  = "win"
	fail = "fail"

	DB *ip2location.DB

	conf *config.Config
	Pipe redis.Pipeliner
	Max  int64

	IPLocationPath string

	ctx = context.Background()
)

func ConnectReids() {
	conf = config.Read()
	Max = int64(conf.Rank.Max)
	println("最大排行数量: ", Max)
	fmt.Println("redis连接参数 :", conf.Rank.RedisUrl, conf.Rank.RedisPassword)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Rank.RedisUrl,
		Password: conf.Rank.RedisPassword,
		DB:       0,
		PoolSize: 1000,
	})
	result, err := Rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis数据库连接错误 :", err)
		return
	}
	println(result)
	// Rdb = rdb.Pipeline()
	// Pipe.ZAdd()
	fmt.Printf("******************redis数据库连接成功******************")
	// iplocationInit()
}

func GetAddr() string {
	return conf.Rank.Port
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
	if count.Val() < Max {
		Rdb.ZRemRangeByRank(ctx, key, Max, -1)
	}
}

// 加入服务器列表
func AddGameServers(key, url string, score float64) {
	Rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: url})
}

func GetTopGameServers(key string) []redis.Z {
	return GetGameServers(key, 0, 0)
}

func DeleGameServer(city, url string) {
	var r = GetSelfCityRank(city, url)
	println("")
	Rdb.ZRemRangeByRank(ctx, city, r-1, r-1)
}

func DeleRank(city, uid string) {
	var r = GetSelfCityRank(city, uid)
	if r < Max {
		Rdb.ZRemRangeByRank(ctx, city, r-1, r-1)
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
	AddRank(worldRank, uid, score)
	AddRank(city, uid, score)
}

// 获取排行榜
func GetRank(key string, start int64, end int64) []redis.Z {
	return Rdb.ZRevRangeWithScores(ctx, key, start, end).Val()
}

func GetGameServers(key string, start int64, end int64) []redis.Z {
	return Rdb.ZRangeWithScores(ctx, key, start, end).Val()
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
	UpdateRank(r.City, r.UID, float64(r.Val))
	return c.JSON(http.StatusOK, r)
}

func RankPull(c echo.Context) error {
	return c.JSON(http.StatusOK, GetWorldRank(0, int64(conf.Rank.WorldRankCount)))
}

func RankCityPull(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	backInfo := GetCityRank(r.City, 0, int64(conf.Rank.CityRankCount))
	return c.JSON(http.StatusOK, backInfo)
}

func RankSheepCityPull(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	backInfo := GetCityRank(r.City, 0, 35)
	return c.JSON(http.StatusOK, backInfo)
}

func RankCityDele(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	DeleRank(r.City, r.UID)
	return c.JSON(http.StatusOK, "")
}

func GetSelf(c echo.Context) error {
	r := &msg.Rank{}
	// 解析body
	ParseNetBody(r, c.Request().Body)
	addr := c.RealIP()
	oldR := &msg.Rank{}
	oldR.WorldRank = GetSelfWorldRank(r.UID)
	oldR.CityRank = GetSelfCityRank(r.City, r.UID)
	oldR.Addr = addr
	return c.JSON(http.StatusOK, oldR)
}

func ParseNetBody(i interface{}, r io.ReadCloser) {
	d, _ := ioutil.ReadAll(r)
	json.Unmarshal(d, i)
}
