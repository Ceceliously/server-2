package user_service

import (
	"server-2/internal/storage"
	"server-2/internal/service/user_service/handlers/v1"
	"server-2/internal/service/user_service/usecase/user"
)

type UserService struct {
	UseCase   *user.UserUseCase
	HandlersV1 *v1.UserHandlersV1
}

func NewUserService(storage storage.UserStorage) *UserService {
	uc := user.NewUserUseCase(storage)
	hV1 := v1.NewHandlerV1(uc)

	return &UserService{
		UseCase:    uc,
		HandlersV1: hV1,
	}
}