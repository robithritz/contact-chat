package main

import (
	"contact-chat/database"
	"contact-chat/middlewares"
	"contact-chat/models/rooms"
	"contact-chat/mqtt"
	chatsRouter "contact-chat/routers/chats"
	"contact-chat/users"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	server.Static("/assets", "static")

	// server.Use(middleware.Logger())

	server.POST("/login", middlewares.Login)
	server.POST("/register", users.Register)

	profileRouters := server.Group("/profile")
	profileRouters.Use(middlewares.AuthorizationCheck)
	profileRouters.GET("", users.GetProfile)

	roomRouters := server.Group("/rooms")
	roomRouters.Use(middlewares.AuthorizationCheck)
	roomRouters.POST("", rooms.CreateRoom)

	chatRouters := server.Group("/chats")
	chatRouters.Use(middlewares.AuthorizationCheck)
	chatRouters.POST("", chatsRouter.SaveMessage)
	chatRouters.GET("/:messageId/received", chatsRouter.ReceivedMessage)
	chatRouters.GET("/:messageId/read", chatsRouter.ReadMessage)

	server.GET("/testmqtt", func(c echo.Context) error {
		mqtt.Publish("haha/test", "oke", 0)
		return c.JSON(http.StatusOK, echo.Map{
			"message": "oke",
		})
	})

	connectedDb, err := database.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(connectedDb)

	mqtt.ConnectBroker()

	server.Logger.Fatal(server.Start(":8888"))
}
