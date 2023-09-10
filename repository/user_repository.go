package repository

import (
	db "cnfs/db/sqlc"
	"cnfs/domain"
)

type userRepository struct {
  db db.Store
}

func NewUserRepository(db db.Store) domain.IUserRepository {
  return &userRepository{db: db}
}

func (r userRepository) Create(user domain.IUser) (domain.IUser, error) {
  // TODO!: update query from db
  return user, nil
}

func (r userRepository) List(page, limit int32) ([]domain.IUser, error) {
  // TODO!: update query from db
  return make([]domain.IUser, limit), nil
}

func (r userRepository) FindByID(id domain.IUserID) (domain.IUser, error) {
  // TODO!: update query from db
  return domain.User{}, nil
}

func (r userRepository) FindByUsername(username domain.IUsername) (domain.IUser, error) {
  // TODO!: update query from db
  return domain.User{}, nil
}
