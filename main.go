package main

import (
	"contact-chat/database"
	"contact-chat/middlewares"
	"contact-chat/models/chats"
	"contact-chat/models/rooms"
	"contact-chat/users"
	"fmt"

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
	chatRouters.POST("", chats.SaveMessage)
	chatRouters.GET("/:messageId/received", chats.ReceivedMessage)
	chatRouters.GET("/:messageId/read", chats.ReadMessage)

	connectedDb, err := database.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(connectedDb)

	server.Logger.Fatal(server.Start(":8888"))
}
