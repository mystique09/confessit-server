package handlers

import (
	"confessit/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var MUser = models.User{}

func GetUsers(conn *gorm.DB) []models.UserResponse {
	var users []models.UserResponse

	conn.Model(&MUser).Scan(&users)
	return users
}

func GetUser(conn *gorm.DB, name string) models.User {
	var user models.User

	conn.Model(&MUser).Where("username = ?", name).Find(&user)
	return user
}

func CreateUser(conn *gorm.DB, payload models.UserCreatePayload) error {
	var nuser models.User = models.User{
		ID:        uuid.New(),
		Username:  payload.Username,
		Password:  payload.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := nuser.HashPassword(); err != nil {
		return err
	}

	if err := conn.Create(&nuser).Error; err != nil {
		return err
	}

	return nil
}

func UpdateUser(conn *gorm.DB, uid uuid.UUID, payload models.UserUpdatePayload) error {
	updatedUser := models.User{
		Username:  payload.Username,
		Password:  payload.Password,
		UpdatedAt: time.Now(),
	}

	if err := updatedUser.HashPassword(); err != nil {
		return err
	}

	if err := conn.Model(&MUser).Where("id = ?", uid).Updates(updatedUser).Error; err != nil {
		return err
	}

	return nil
}

func DeleteUser(conn *gorm.DB, uid uuid.UUID) error {
	if err := conn.Delete(&MUser, "id = ?", uid).Error; err != nil {
		return err
	}

	return nil
}
