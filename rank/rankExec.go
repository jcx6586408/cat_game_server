package main

import (
	"leafserver/src/server/conf"
	"net/http"
	"rank/rank"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:          nil,
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	rank.RankInit()
	e.POST("/self", rank.GetSelf)
	e.POST("/pull", rank.RankPull)
	e.POST("/cityPull", rank.RankCityPull)
	e.POST("/update", rank.RankUpdate)
	e.POST("/roomCreate", rank.RoomCreate)
	e.POST("/openid", rank.GetOpenID)

	sysType := runtime.GOOS

	if sysType == "linux" {
		// LINUX系统
		e.Logger.Fatal(e.StartTLS(conf.Server.HttpAddr, conf.Server.CertFile, conf.Server.KeyFile))
	}

	if sysType == "windows" {
		// windows系统
		e.Logger.Fatal(e.Start(conf.Server.HttpAddr))
	}

}
