package main

import (
	"fmt"
	"leafserver/src/server/conf"
	"net/http"
	"os"
	"os/signal"
	"rank/rank"
	"runtime"
	"syscall"

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
		AllowHeaders:     []string{"*"},
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

	OnExit()

	if sysType == "linux" {
		// LINUX系统
		e.Logger.Fatal(e.StartTLS(conf.Server.HttpAddr, conf.Server.CertFile, conf.Server.KeyFile))
	}

	if sysType == "windows" {
		// windows系统
		e.Logger.Fatal(e.Start(conf.Server.HttpAddr))
	}

}

func OnExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		for {
			s := <-ch
			switch s {
			case syscall.SIGINT:
				//SIGINT 信号，在程序关闭时会收到这个信号
				fmt.Println("SIGINT", "退出程序，执行退出前逻辑")
				return
			case syscall.SIGKILL:
				fmt.Println("SIGKILL")
				return
			default:
				fmt.Println("default")
			}
		}
	}()

}
