package main

import (
	"config"
	"excel"
	"leafserver/src/server/conf"
	"net/http"
	"os"
	"rank/rank"
	"runtime"
	"storage/redis"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	conf.ConfPath = os.Args[1]
	config.ConfPath = os.Args[2]
	config.RoomConfPath = os.Args[3]
	excel.TablePath = os.Args[4]
	redis.IPLocationPath = os.Args[5]
	if len(os.Args) >= 7 {
		conf.Server.CertFile = os.Args[6]
	}
	if len(os.Args) >= 8 {
		conf.Server.KeyFile = os.Args[7]
	}
	redis.ConnectReids()
	e := echo.New()
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:          nil,
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	rank.RankInit()
	e.POST("/self", redis.GetSelf)
	e.POST("/pull", redis.RankPull)
	e.POST("/cityPull", redis.RankCityPull)
	e.POST("/update", redis.RankUpdate)

	sysType := runtime.GOOS

	if sysType == "linux" {
		// LINUX系统
		e.Logger.Fatal(e.StartTLS(redis.GetAddr(), conf.Server.CertFile, conf.Server.KeyFile))
	}

	if sysType == "windows" {
		// windows系统
		e.Logger.Fatal(e.Start(redis.GetAddr()))
	}
}
