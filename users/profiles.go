package users

import (
	"contact-chat/database"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Profile struct {
	User
	CreatedAt time.Time `json:"createdAt"`
}

func GetProfile(c echo.Context) error {
	profile := &Profile{}
	username := c.Get("username")

	result := database.DB.Model(&database.Users{}).Debug().Select("username, first_name, last_name, image_url, phone_number, created_at").First(&profile, "username = ?", username)
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{})
	}

	return c.JSON(http.StatusOK, profile)

}
