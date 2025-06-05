package v1

import (
	"net/http"
	"errors"
	"log"
	"github.com/go-chi/render"
	"server-2/internal/storage"
	"server-2/internal/service/user_service/usecase/user"
	usc "server-2/internal/models/user/user_create"
	ug "server-2/internal/models/user/user_get"
	res "server-2/internal/models/response"
	"github.com/go-playground/validator/v10"
)

type UserHandlersV1 struct {
	uc *user.UserUseCase
}

func NewHandlerV1(uc *user.UserUseCase) *UserHandlersV1 {
	return &UserHandlersV1{uc:uc}
}


func (h *UserHandlersV1) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
		var req usc.UserCreateRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			render.JSON(w, r, res.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			render.JSON(w, r, res.Error("invalid request"))
			return
		}


		err = h.uc.CreateUser(req)
		if errors.Is(err, storage.ErrUserExists) {
			render.JSON(w, r, res.Error("user already exists"))
			return
		}
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				render.JSON(w, r, res.Error("failed to add user"))
				return
			}

	log.Printf("user %s is added", req.Username)

	render.JSON(w, r, usc.ToResponse(req.Username))

	}

	func (h *UserHandlersV1) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	username, ok := r.Context().Value("username").(string)
		if !ok {
			render.JSON(w, r, res.Error("internal server error"))
			return
		}

	ObtainedUser, err := h.uc.GetUser(username)
	if errors.Is(err, storage.ErrUserNotFound) {
		render.JSON(w, r, res.Error("failed to find user"))
		return
	}
	if err != nil {
			render.JSON(w, r, res.Error("failed to find user"))
			return
		}

	log.Printf("user %s is found", username)

	render.JSON(w, r, ug.ToResponse(ObtainedUser))
	}