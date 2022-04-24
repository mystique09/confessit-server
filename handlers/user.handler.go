package handlers

import (
	"confessit/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUsers(conn *gorm.DB) []models.UserResponse {
	return make([]models.UserResponse, 2)
}

func GetUser(conn *gorm.DB, name string) models.UserResponse {
	return models.UserResponse{}
}

func CreateUser(conn *gorm.DB, payload models.UserCreatePayload) models.UserResponse {
	return models.UserResponse{}
}

func UpdateUser(conn *gorm.DB, uid uuid.UUID, payload models.UserUpdatePayload) bool {
	return false
}

func DeleteUser(conn *gorm.DB, uid uuid.UUID) bool {
	return false
}
