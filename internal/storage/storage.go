package storage

import (
	"errors"
	entity "server-2/server/service/handlers/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user exists")
)



type UserStorage interface {
	Create(user *entity.User) (error)
	GetUser(username string) (*entity.User, error)
}