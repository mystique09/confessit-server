package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken error = errors.New("token has expired")
var ErrInvalidToken error = errors.New("token is invalid")

type Payload struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userId uuid.UUID, username string, duration time.Duration) (*Payload, error) {
	tokenId := uuid.New()
	payload := &Payload{
		Id:        tokenId,
		UserId:    userId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
