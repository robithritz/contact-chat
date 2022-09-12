package chats

import (
	"contact-chat/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SaveMessage(c echo.Context) error {
	message := new(database.TChats)

	if err := c.Bind(message); err != nil {
		return echo.ErrBadRequest
	}

	senderId, err := strconv.Atoi(fmt.Sprintf("%v", c.Get("userId")))
	if err != nil {
		return echo.ErrBadRequest
	}

	message.SenderID = uint(senderId)

	result := database.DB.Create(&message)
	if result.Error != nil {
		return echo.ErrBadGateway
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "message created",
	})
}

func ReceivedMessage(c echo.Context) error {
	messageId, err := strconv.Atoi(c.Param("messageId"))
	if err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.Atoi(fmt.Sprintf("%v", c.Get("userId")))
	if err != nil {
		return echo.ErrBadRequest
	}

	tChatSent := new(database.TChatSents)
	recordExisted := database.DB.Where("chat_id = ? AND target_id = ?", messageId, userId).First(&tChatSent)
	if recordExisted.RowsAffected == 0 {
		tChatSent.ChatID = uint(messageId)
		tChatSent.TargetID = uint(userId)
		result := database.DB.Create(&tChatSent)
		if result.Error != nil {
			return echo.ErrBadGateway
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "status updated successfully",
	})
}

func ReadMessage(c echo.Context) error {
	messageId, err := strconv.Atoi(c.Param("messageId"))
	if err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.Atoi(fmt.Sprintf("%v", c.Get("userId")))
	if err != nil {
		return echo.ErrBadRequest
	}

	tChatReaders := new(database.TChatReaders)
	recordExisted := database.DB.Where("chat_id = ? AND target_id = ?", messageId, userId).First(&tChatReaders)
	if recordExisted.RowsAffected == 0 {
		tChatReaders.ChatID = uint(messageId)
		tChatReaders.TargetID = uint(userId)
		result := database.DB.Create(&tChatReaders)
		if result.Error != nil {
			return echo.ErrBadGateway
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "status updated successfully",
	})
}
