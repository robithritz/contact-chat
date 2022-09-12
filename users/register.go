package users

import (
	"contact-chat/database"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	ImageURL    string `json:"imageURL"`
	PhoneNumber string `json:"phoneNumber"`
}

type RegisterUser struct {
	User
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	userRegister := new(RegisterUser)
	if err := c.Bind(userRegister); err != nil {
		return echo.ErrBadRequest
	}

	var userExisted User
	result := database.DB.Model(&database.Users{}).First(&userExisted, "username = ?", userRegister.Username)
	if result.RowsAffected > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "username already existed",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 12)
	if err != nil {
		return echo.ErrBadGateway
	}

	submitForm := &database.Users{
		Username:    userRegister.Username,
		Password:    string(password),
		FirstName:   userRegister.FirstName,
		LastName:    userRegister.LastName,
		ImageURL:    userRegister.ImageURL,
		PhoneNumber: userRegister.PhoneNumber,
	}

	created := database.DB.Create(&submitForm)
	return c.JSON(http.StatusOK, created.RowsAffected)

}
