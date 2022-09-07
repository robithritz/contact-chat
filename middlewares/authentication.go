package middlewares

import (
	"contact-chat/database"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type jwtClaims struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	ImageURL    string `json:"imageURL"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.StandardClaims
}

type UserInfo struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	ImageURL    string `json:"imageURL"`
	PhoneNumber string `json:"phoneNumber"`
}

type returnMessage struct {
	Message string `json:"message"`
}

var secrets = "404142434445464748494A4B4C4D4E4F"

var AuthCheck = &middleware.JWTConfig{
	SigningKey: []byte(secrets),
	Claims:     &jwtClaims{},
}

func Login(c echo.Context) error {
	userLogin := new(UserLogin)
	if err := c.Bind(userLogin); err != nil {
		return echo.ErrBadRequest
	}

	if userLogin.Username == "" || userLogin.Password == "" {
		return c.JSON(http.StatusBadRequest, &returnMessage{
			Message: "username and password cannot be empty",
		})
	}

	users := &database.Users{}

	result := database.DB.First(&users, "username = ?", userLogin.Username)
	fmt.Println(result.Error)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return c.JSON(http.StatusBadGateway, &returnMessage{
			Message: "something went wrong!",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(userLogin.Password)); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, &returnMessage{
			Message: "wrong username or password",
		})
	}

	userClaim := &jwtClaims{
		Username:    users.Username,
		FirstName:   users.FirstName,
		LastName:    users.LastName,
		ImageURL:    users.ImageURL,
		PhoneNumber: users.PhoneNumber,
	}

	dataSigning := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)

	token, err := dataSigning.SignedString([]byte(secrets))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":   token,
		"message": "successful",
	})
}

func AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			c.Error(echo.ErrUnauthorized)
			return nil
		}

		token, _ := jwt.ParseWithClaims(authHeader, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secrets), nil
		})

		if claims, ok := token.Claims.(*jwtClaims); !ok || !token.Valid {
			c.Error(echo.ErrUnauthorized)
			return nil
		} else {
			/* Check if user still active */
			profile := &database.Users{}
			userExisted := database.DB.First(&profile, "username = ?", claims.Username)
			if userExisted.RowsAffected == 0 {
				c.Error(echo.ErrUnauthorized)
				return nil
			}

			c.Set("username", claims.Username)
			return next(c)
		}
	}
}
