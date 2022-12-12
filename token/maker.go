package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(userId uuid.UUID, username string, duration time.Duration) (uuid.UUID, string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
