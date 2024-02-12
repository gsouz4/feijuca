package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("", nil)
	e.POST("", nil)

	e.Logger.Fatal(e.Start(":8080"))
}
