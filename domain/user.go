package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	IUserID interface {
		IBaseField
		From(id uuid.UUID) IUserID
	}
	IUsername = IBaseField
	IPassword = IBaseField

	UserID struct {
		value uuid.UUID
	}

	Username struct {
		value string
	}

	Password struct {
		value string
	}
)

func NewUserID() IUserID {
	return UserID{value: uuid.New()}
}

func (id UserID) String() string {
	return id.value.String()
}

func (UserID) From(id uuid.UUID) IUserID {
	return UserID{value: id}
}

func (id UserID) ValidateLength(n int) bool {
	return len(id.String()) == n
}

func NewUsername(value string) IUsername {
	return Username{value: value}
}

func (u Username) ValidateLength(n int) bool {
	return len(u.value) >= n
}

func (u Username) String() string {
	return u.value
}

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
	IUser interface {
		ID() IUserID
		Username() IUsername
		Password() IPassword
		IDateFields
	}

	User struct {
		id         IUserID
		username   IUsername
		password   IPassword
		created_at time.Time
		updated_at time.Time
	}
)

func NewUser(username IUsername, password IPassword) IUser {
	return User{
		id:         NewUserID(),
		username:   username,
		password:   password,
		created_at: time.Now(),
		updated_at: time.Now(),
	}
}

func (u User) ID() IUserID {
	return u.id
}

func (u User) Username() IUsername {
	return u.username
}

func (u User) Password() IPassword {
	return u.password
}

func (u User) CreatedAt() time.Time {
	return u.created_at
}

func (u User) UpdatedAt() time.Time {
	return u.updated_at
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) IntoResponse() Response[UserResponse] {
	return Response[UserResponse]{
		Message: "",
		Data: UserResponse{
			ID:        u.ID().String(),
			Username:  u.Username().String(),
			CreatedAt: u.CreatedAt(),
			UpdatedAt: u.UpdatedAt(),
		},
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
