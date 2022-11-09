package rooms

import (
	"contact-chat/database"
	"contact-chat/models/chats"
	"contact-chat/mqtt"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoomFull struct {
	database.MRooms
	Participants []uint `json:"participants"`
}

func CreateRoom(c echo.Context) error {
	newRoomFull := new(RoomFull)

	if err := c.Bind(newRoomFull); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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

		/*
			publish to all participant of created room to inform them to subscribe to room's topic, except the creator
		*/
		if val != uint(userId) {
			partId := strconv.Itoa(int(val))
			newMQID := uuid.New()
			payload := echo.Map{
				"roomId":          newRoom.ID,
				"creatorId":       userId,
				"roomName":        newRoomFull.RoomName,
				"roomType":        newRoomFull.RoomType,
				"roomDescription": newRoomFull.RoomDescription,
				"roomImageURL":    newRoomFull.RoomImageURL,
				"participants":    newRoomFull.Participants,
			}

			payloadStringified, _ := json.Marshal(payload)
			targetTopic := "users/" + partId + "/new-room/" + newMQID.String()

			logMQ := &database.LMQ{
				MQID:        newMQID,
				TargetTopic: targetTopic,
				TargetID:    val,
				Payload:     string(payloadStringified),
			}
			err := chats.SaveMQLog(logMQ)
			if err != nil {
				return c.JSON(http.StatusBadGateway, echo.Map{
					"message": "Something went wrong",
				})
			}

			mqtt.Publish(targetTopic, payload, 1)
		}
	}

	database.DB.Create(&listParticipants)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "room successfully created",
		"roomId":  newRoom.ID,
	})
}
