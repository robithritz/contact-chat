package rooms

import (
	"contact-chat/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoomFull struct {
	database.MRooms
	Participants []uint `json:"participants"`
}

func CreateRoom(c echo.Context) error {
	newRoomFull := new(RoomFull)

	if err := c.Bind(newRoomFull); err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.Atoi(fmt.Sprintf("%v", c.Get("userId")))
	if err != nil {
		return echo.ErrBadGateway
	}
	newRoomFull.CreatorID = uint(userId)
	newRoom := &database.MRooms{
		RoomName:        newRoomFull.RoomName,
		RoomType:        newRoomFull.RoomType,
		RoomImageURL:    newRoomFull.RoomImageURL,
		RoomDescription: newRoomFull.RoomDescription,
		CreatorID:       uint(userId),
	}

	result := database.DB.Create(&newRoom)
	if result.Error != nil {
		return echo.ErrBadGateway
	}

	var listParticipants []database.MRoomParticipants

	for _, val := range newRoomFull.Participants {
		listParticipants = append(listParticipants, database.MRoomParticipants{
			RoomId:        newRoom.ID,
			ParticipantId: val,
		})
	}

	database.DB.Create(&listParticipants)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "room successfully created",
		"roomId":  newRoom.ID,
	})
}
