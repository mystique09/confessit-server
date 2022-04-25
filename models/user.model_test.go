package models

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate = validator.New()

func TestUserHashPasswordAlgo(t *testing.T) {
	nuser := User{
		ID:        uuid.New(),
		Username:  "testuser133",
		Password:  "testpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	raw_pass := nuser.Password

	if assert.NoError(t, nuser.HashPassword()) {
		assert.Equal(t, nuser.Username, "testuser133")
		assert.NotEqual(t, raw_pass, nuser.Password)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(nuser.Password), []byte("testpassword")))
	}
}

func TestUserValidation(t *testing.T) {
	ncuser := UserCreatePayload{
		Username: "te",
		Password: "tes",
	}

	assert.Error(t, validate.Struct(ncuser))
}

func TestUserValidationSuccess(t *testing.T) {
	ncuser := UserCreatePayload{
		Username: "testuser",
		Password: "testpasword",
	}

	assert.NoError(t, validate.Struct(ncuser))
}
