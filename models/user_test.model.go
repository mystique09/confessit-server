package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TesUserHashPasswordAlgo(t *testing.T) {
	nuser := User{
		ID:        uuid.New(),
		Username:  "Test",
		Password:  "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	raw_pass := nuser.Password

	if assert.NoError(t, nuser.HashPassword()) {
		assert.Equal(t, nuser.Username, "Test")
		assert.NotEqual(t, raw_pass, nuser.Password)
	}
}
