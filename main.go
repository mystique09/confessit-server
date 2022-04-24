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

	route := routers.Route{Conn: conn}

	app := echo.New()
	app.Use(routers.CustomLoggerMiddleware())

	app.GET("/", routers.ConfessIt)
	app.POST("/auth", route.Login)
	app.POST("/signup", route.Signup)

	user_g := app.Group("/users")
	{
		user_g.GET("", route.GetUsers)
		user_g.GET("/:name", route.GetUser)
		user_g.PUT("/:id", route.UpdateUser)
		user_g.DELETE("/:id", route.DeleteUser)
	}

	message_g := app.Group("/messages")
	{
		message_g.GET("", route.GetMessages)
		message_g.POST("", route.CreateMessage)
		message_g.GET("/:name", route.GetMessage)
		message_g.DELETE("/:id", route.DeleteMessage)
	}

	app.Logger.Fatal(app.Start(":5000"))
}
