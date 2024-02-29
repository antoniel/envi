package main

import (
	"envii/apps/api/initializers"

	"github.com/labstack/echo"
)

func init() {
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	e := echo.New()
	RegisterMiddlewares(e)
	RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
