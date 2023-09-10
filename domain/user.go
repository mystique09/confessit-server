package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
  IUsername interface {
    ValidateLength(n int) bool
    String() string
  }

  Username struct {
    value string
  }
)

func NewUsername(value string) IUsername {
  return Username{value: value}
}

func (u Username) ValidateLength(n int) bool {
  return len(u.value) >= n
}

func (u Username) String() string {
  return u.value
}

type (
  IPassword interface {
    ValidateLength(n int) bool
    String() string
  }

  Password struct {
    value string
  }
)

func NewPassword(value string) IPassword {
  return Password{value: value}
}

func (p Password) ValidateLength(n int) bool {
  return len(p.value) >= n
}

func (p Password) String() string {
  return p.value
}

type (
  IUserID interface {
    String() string
  }

  UserID struct {
    value uuid.UUID
  }
)

func NewUserID() IUserID {
  return UserID{value: uuid.New()}
}

func (id UserID) String() string {
  return id.value.String()
}

type (
  IUser interface {
    ID() string
    Username() string
    Password() string
    CreatedAt() time.Time
    UpdatedAt() time.Time
  }

  User struct {
    id IUserID
    username IUsername
    password IPassword
    created_at time.Time
    updated_at time.Time
  }
)

func NewUser(username IUsername, password IPassword) IUser {
  return User{
    id: NewUserID(),
    username: username,
    password: password,
  }
}

func (u User) ID() string {
  return u.id.String()
}

func (u User) Username() string {
  return u.username.String()
}

func (u User) Password() string {
  return u.password.String()
}

func (u User) CreatedAt() time.Time {
  return u.created_at
}

func (u User) UpdatedAt() time.Time {
  return u.updated_at
}

type UserResponse struct {
  ID string `json:"id"`
  Username string `json:"username"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func (u User) IntoResponse() UserResponse {
  return UserResponse{
    ID: u.ID(),
    Username: u.Username(),
    CreatedAt: u.CreatedAt(),
    UpdatedAt: u.UpdatedAt(),
  }
}

type CreateUserDTO struct {
  Username string `json:"username" validate:"required,min=1,max=20"`
  Password string `json:"password" validate:"required,min=8,max=20"`
}

func (payload CreateUserDTO) ToUser() IUser {
  return NewUser(NewUsername(payload.Username), NewPassword(payload.Password))
}

type IUserRepository interface {
  Create(user IUser) (IUser, error)
  List(page, limit int32) ([]IUser, error)
  FindByID(id IUserID) (IUser, error)
  FindByUsername(username IUsername) (IUser, error)
}

type ISignupUserUseCase interface {
  CheckUsernameAvailability(username Username) error
  Signup(payload CreateUserDTO) (UserResponse, error)
}
