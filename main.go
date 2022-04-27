package main

import (
	"confessit/db"
	"confessit/models"
	"confessit/routers"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	conn := db.ConnectDb()
	conn.AutoMigrate(models.User{}, models.Message{})

	route := routers.Route{Conn: conn}

	app := echo.New()
	app.Use(routers.CustomLoggerMiddleware())
	app.Use(routers.CustomCORSMiddleware())
	//app.Use(routers.CustomCSRFMiddleware())
	app.Use(routers.CustomRateLimitMiddleware())

	app.GET("/", routers.ConfessIt)
	app.POST("/auth", route.Login)
	app.POST("/signup", route.Signup)

	user_g := app.Group("/users", routers.AuthMiddleware())
	{
		user_g.GET("", route.GetUsers)
		user_g.GET("/:name", route.GetUser)
		user_g.PUT("/:name", route.UpdateUser)
		user_g.DELETE("/:name", route.DeleteUser)
	}

	message_g := app.Group("/messages", routers.AuthMiddleware())
	{
		message_g.GET("", route.GetMessages)
		message_g.POST("", route.CreateMessage)
		message_g.GET("/:name", route.GetMessage)
		message_g.DELETE("/:id", route.DeleteMessage)
	}

	var port string = os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	app.Logger.Fatal(app.Start(":" + port))
}
