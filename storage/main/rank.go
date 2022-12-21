package main

import (
	"leafserver/src/server/conf"
	"net/http"
	"rank/rank"
	"runtime"
	"storage/redis"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
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
