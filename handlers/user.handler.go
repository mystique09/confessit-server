package handlers

import (
	"confessit/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var MUser = models.User{}

func GetUsers(conn *gorm.DB) []models.UserResponse {
	var users []models.UserResponse

	conn.Model(&MUser).Scan(&users)
	return users
}

func GetUser(conn *gorm.DB, name string) models.UserResponse {
	var user models.UserResponse

	conn.Model(&MUser).Where("username = ?", name).Find(&user)
	return user
}

func CreateUser(conn *gorm.DB, payload models.UserCreatePayload) error {
	hash_pass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	var nuser models.User = models.User{
		ID:        uuid.New(),
		Username:  payload.Username,
		Password:  string(hash_pass),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := conn.Create(&nuser).Error; err != nil {
		return err
	}

	return nil
}

func UpdateUser(conn *gorm.DB, uid uuid.UUID, payload models.UserUpdatePayload) error {
	if err := conn.Model(&MUser).Where("id = ?", uid).Updates(
		models.User{
			Username:  payload.Username,
			Password:  payload.Password,
			UpdatedAt: time.Now(),
		},
	).Error; err != nil {
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
