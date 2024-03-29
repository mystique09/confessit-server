// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Satisfaction string

const (
	SatisfactionLIKE    Satisfaction = "LIKE"
	SatisfactionDISLIKE Satisfaction = "DISLIKE"
)

func (e *Satisfaction) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Satisfaction(s)
	case string:
		*e = Satisfaction(s)
	default:
		return fmt.Errorf("unsupported scan type for Satisfaction: %T", src)
	}
	return nil
}

type NullSatisfaction struct {
	Satisfaction Satisfaction
	Valid        bool // Valid is true if Satisfaction is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSatisfaction) Scan(value interface{}) error {
	if value == nil {
		ns.Satisfaction, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Satisfaction.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSatisfaction) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Satisfaction, nil
}

type Comment struct {
	ID             uuid.UUID `json:"id"`
	Content        string    `json:"content"`
	UserIdentityID uuid.UUID `json:"user_identity_id"`
	PostID         uuid.UUID `json:"post_id"`
	ParentID       uuid.UUID `json:"parent_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Message struct {
	ID         uuid.UUID `json:"id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	Seen       bool      `json:"seen"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Post struct {
	ID             uuid.UUID `json:"id"`
	Content        string    `json:"content"`
	UserIdentityID uuid.UUID `json:"user_identity_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserIdentity struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	IdentityHash uuid.UUID `json:"identity_hash"`
}
