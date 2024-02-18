package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, remote_ip=${remote_ip}, method=${method}, " +
			"path=${path}, status=${status}, took=${response_time}, sent=t=${response_size} bytes\n",
	}))
}
func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", healthHandler)
	e.GET("/oauth/authorize/google", oAuthAuthorizeGoogleHandler)
	e.GET("/oauth/callback/google", oAuthCallbackGoogleHandler)
}
