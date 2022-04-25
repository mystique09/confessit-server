package models

import "time"

type Message struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

type MessageCreatePayload struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required,min=1,max=1500"`
}

type MessageDeletePayload struct {
	ID int `json:"id" validate:"required,numeric"`
}
