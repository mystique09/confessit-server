package handlers

import (
	"confessit/models"
	"time"

	"gorm.io/gorm"
)

var MMessage = models.Message{}

func GetMessages(conn *gorm.DB, uname string) []models.Message {
	var messages []models.Message

	conn.Model(&MMessage).Where("to = ?", uname).Find(&messages)

	return messages
}

func GetMessage(conn *gorm.DB, uname string) models.Message {
	var message models.Message

	conn.Model(&MMessage).Where("to = ?", uname).Find(&message)

	return message
}

func CreateMessage(conn *gorm.DB, payload models.MessageCreatePayload) error {
	var total_message = GetMessages(conn, payload.To)
	var nmessage models.Message = models.Message{
		ID:        uint32(len(total_message) + 1),
		To:        payload.To,
		Message:   payload.Message,
		CreatedAt: time.Now(),
	}

	if err := conn.Create(&nmessage).Error; err != nil {
		return err
	}

	return nil
}

func DeleteMessage(conn *gorm.DB, uid uint32) error {
	if err := conn.Delete(&MMessage, "id = ?", uid).Error; err != nil {
		return err
	}

	return nil
}
