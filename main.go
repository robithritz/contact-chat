package main

import (
	"contact-chat/database"
	"contact-chat/middlewares"
	"contact-chat/profiles"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	server.Static("/assets", "static")

	// server.Use(middleware.Logger())

	server.POST("/login", middlewares.Login)

	profileRouters := server.Group("/profile")
	profileRouters.Use(middlewares.AuthorizationCheck)

	profileRouters.GET("", profiles.GetProfile)

	connectedDb, err := database.Connect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(connectedDb)

	server.Logger.Fatal(server.Start(":8888"))
}

func getHandler(c echo.Context) error {
	name := c.Param("name")
	age := c.QueryParam("age")
	return c.String(http.StatusOK, "Hello, "+name+"! you are "+age+" years old\n")
}

func postHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, "oke")
}

func basicAuth(username string, password string, c echo.Context) (bool, error) {
	if username == "admin" && password == "admin" {
		return true, nil
	}
	return false, nil
}
