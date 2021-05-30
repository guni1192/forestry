package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(logging bool) *echo.Echo {
	e := echo.New()

	if logging {
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
	}

	e.GET("/health", Health)

	return e
}
