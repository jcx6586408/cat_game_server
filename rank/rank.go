package main

import (
	"net/http"
	"rank/rank"

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
	e.Logger.Fatal(e.Start(":8443"))
}
