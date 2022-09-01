package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Person struct {
	Name string `json:"name" form:"name" query:"name"`
	Age  int    `json:"age"`
}

func main() {
	server := echo.New()
	server.GET("/:name", getHandler)
	server.POST("/", postHandler)

	server.Static("/assets", "static")

	// server.Use(middleware.Logger())

	adminRouter := server.Group(("/admin"))
	adminRouter.Use(middleware.BasicAuth(basicAuth))

	adminRouter.GET("/data", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "okee")
	})

	connectedDb, err := connect()
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
	person := new(Person)
	// name := c.FormValue("name")
	if err := c.Bind(person); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, person)
}

func basicAuth(username string, password string, c echo.Context) (bool, error) {
	if username == "admin" && password == "admin" {
		return true, nil
	}
	return false, nil
}
