package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user exists")
)



type UserStorage interface {
	Create(username, password string, firstName, lastName *string, age *int) (error)
	GetUser(username, password string) (*string, *string, *int, error)
}