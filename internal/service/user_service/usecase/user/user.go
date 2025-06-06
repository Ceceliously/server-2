package user

import (
	"fmt"
	"errors"

	"server-2/internal/storage"
	serv "server-2/internal/models/user/user_create"
	"golang.org/x/crypto/bcrypt"
	entity "server-2/internal/models/user/user"
)


type UserUseCase struct {
	storage storage.UserStorage
}

func NewUserUseCase (s storage.UserStorage) *UserUseCase {
	return &UserUseCase{storage: s}
}

func (uc *UserUseCase) CreateUser(req serv.UserCreateRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash pass: %w", err)
	}

	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		FirstName: req.FirstName,
		LastName: req.LastName,
		Age: req.Age,
	}

	err = uc.storage.Create(user)
		if errors.Is(err, storage.ErrUserExists) {
			return storage.ErrUserExists
		}
		return err
}

func (uc *UserUseCase)GetUser(username string) (*entity.User, error) {
	msg, err :=  uc.storage.GetUser(username)
	if errors.Is(err, storage.ErrUserNotFound) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
			return nil, err
		}
	return msg, nil
}


