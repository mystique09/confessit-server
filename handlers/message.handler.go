package handlers

import (
	"confessit/models"

	"gorm.io/gorm"
)

func GetMessages(conn *gorm.DB) []models.Message {
	return make([]models.Message, 2)
}

func GetMessage(conn *gorm.DB, name string) models.Message {
	return models.Message{}
}

func CreateMessage(conn *gorm.DB, payload models.MessageCreatePayload) models.Message {
	return models.Message{}
}

func UpdateMessage(conn *gorm.DB, uid uint32, payload models.UserUpdatePayload) bool {
	return false
}

func DeleteMessage(conn *gorm.DB, uid uint32) bool {
	return false
}
