package domain

import "fmt"

var (
  MAX_USERNAME_LENGTH = 20
  MAX_PASSWORD_LENGTH = 20
  INVALID_USERNAME =   fmt.Sprintf("username must be between 1 and %d characters", MAX_USERNAME_LENGTH)
  INVALID_PASSWORD =   fmt.Sprintf("password must be between 1 and %d characters", MAX_PASSWORD_LENGTH)
  INVALID_CREDENTIALS = "invalid username or password"
  INVALID_EMAIL = "invalid email address"

  USER_NOT_FOUND = "user not found"
  USER_ALREADY_EXISTS = "user already taken"
  USER_NOT_AUTHENTICATED = "user not authenticated"
  EMAIL_ALREADY_EXISTS = "email already taken"
)
