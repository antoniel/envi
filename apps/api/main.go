package main

import (
	"envii/apps/api/initializers"
	libEcho "envii/apps/api/lib/echo"

	"github.com/labstack/echo"
)

func init() {
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	e := echo.New()
	libEcho.RegisterMiddlewares(e)
	libEcho.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
