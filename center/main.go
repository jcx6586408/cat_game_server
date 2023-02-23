package main

import (
	"config"
	"excel"
	"leafserver/src/server/conf"
	"net/http"
	"os"
	"rank/rank"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	println("启动中心服")
	conf.ConfPath = os.Args[1]
	conf.Read()
	config.ConfPath = os.Args[2]
	config.RoomConfPath = os.Args[3]
	excel.TablePath = os.Args[4]
	rank.Port = os.Args[5]
	port := os.Args[6]
	if len(os.Args) >= 8 {
		conf.Server.CertFile = os.Args[7]
	}
	if len(os.Args) >= 9 {
		conf.Server.KeyFile = os.Args[8]
	}
	e := echo.New()
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:          nil,
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"*"},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	rank.ConfInit()
	go func() {
		rank.CenterInit()
	}()
	e.POST("/roomCreate", rank.RoomCreate)
	e.POST("/openid", rank.GetOpenID)
	e.POST("/bytedanceopenid", rank.GetBytedanceOpenID)

	sysType := runtime.GOOS

	if sysType == "linux" {
		// LINUX系统
		e.Logger.Fatal(e.StartTLS(port, conf.Server.CertFile, conf.Server.KeyFile))
	}

	if sysType == "windows" {
		// windows系统
		e.Logger.Fatal(e.Start(port))
	}

}
