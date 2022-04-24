package main

import (
	"confessit/db"
	"confessit/models"
	"confessit/routers"

	"github.com/labstack/echo/v4"
)

func main() {

	conn := db.ConnectDb()
	conn.AutoMigrate(models.User{}, models.Message{})

	app := echo.New()
	app.Use(routers.CustomLoggerMiddleware())

	app.GET("/", routers.ConfessIt)

	app.Logger.Fatal(app.Start(":5000"))
}
