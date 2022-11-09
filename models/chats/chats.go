package chats

import (
	"contact-chat/database"
)

func SaveMessage(data *database.TChats) error {
	result := database.DB.Create(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SaveMQLog(data *database.LMQ) error {
	result := database.DB.Create(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SaveMQLogBulk(data []database.LMQ) error {
	result := database.DB.Create(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ReceivedMessage(data *database.TChatSents) error {
	tChatSent := new(database.TChatSents)
	recordExisted := database.DB.Where("message_id = ? AND target_id = ?", data.MessageId, data.TargetID).First(&tChatSent)
	if recordExisted.RowsAffected == 0 {
		tChatSent.MessageId = data.MessageId
		tChatSent.TargetID = data.TargetID
		tChatSent.SentAt = data.SentAt
		result := database.DB.Create(&tChatSent)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func ReadMessage(data *database.TChatReaders) error {
	tChatRead := new(database.TChatReaders)
	recordExisted := database.DB.Where("message_id = ? AND target_id = ?", data.MessageId, data.TargetID).First(&tChatRead)
	if recordExisted.RowsAffected == 0 {
		tChatRead.MessageId = data.MessageId
		tChatRead.TargetID = data.TargetID
		tChatRead.ReadAt = data.ReadAt
		result := database.DB.Create(&tChatRead)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
