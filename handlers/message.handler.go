package handlers

import (
	"confessit/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var MMessage = models.Message{}

func GetMessages(conn *gorm.DB, uname string) []models.Message {
	var messages []models.Message

	conn.Model(&MMessage).Where(`"to" = ?`, uname).Scan(&messages)

	return messages
}

func GetMessage(conn *gorm.DB, uname string) models.Message {
	var message models.Message

	conn.Model(&MMessage).Where("to = ?", uname).Find(&message)

	return message
}

func CreateMessage(conn *gorm.DB, payload models.MessageCreatePayload) error {
	var nmessage models.Message = models.Message{
		ID:        uuid.New(),
		To:        payload.To,
		Message:   payload.Message,
		CreatedAt: time.Now(),
	}

	if err := conn.Create(&nmessage).Error; err != nil {
		return err
	}

	return nil
}

func DeleteMessage(conn *gorm.DB, uid uuid.UUID, to string) error {
	if err := conn.Delete(&MMessage, `id = ? AND "to" = ?`, uid, to).Error; err != nil {
		return err
	}

	return nil
}
