package domain

import "time"

type (
	IBaseField interface {
		ValidateLength(n int) bool
		String() string
	}

	IDateFields interface {
		CreatedAt() time.Time
		UpdatedAt() time.Time
	}

	Response[T any] struct {
		Message string `json:"message,omitempty"`
		Data    T      `json:"data"`
	}
)
