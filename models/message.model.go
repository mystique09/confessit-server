package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

type MessagePayload struct {
	To string `json:"to" validate:"required"`
}

type MessageCreatePayload struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required,min=1,max=1500"`
}

type MessageDeletePayload struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
