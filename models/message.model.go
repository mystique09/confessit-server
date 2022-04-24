package models

import "time"

type Message struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

type MessageCreatePayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type MessageDeletePayload struct {
	ID int `json:"id"`
}
