package create

import (
	"log"
	"net/http"
	"errors"
	"server-2/internal/storage"
	res "server-2/internal/models/response"
	serv "server-2/internal/models/user/user_create"
	"server-2/internal/service/user_service/usecase/user"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)



func CreateUserHandler(uc *user.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.create.new"

		var req serv.UserCreateRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Printf("%s: failed to decode request body", fn)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			render.JSON(w, r, res.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Printf("invalid request : %w", err)
			render.JSON(w, r, res.Error("invalid request"))
			return
		}


		err = uc.CreateUser(req)
		if errors.Is(err, storage.ErrUserExists) {
			log.Printf("user already exists : %w", err)
			render.JSON(w, r, res.Error("user already exists"))
			return
		}
		if err != nil {
				log.Printf("%s: failed to add user", fn)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				render.JSON(w, r, res.Error("failed to add user"))
				return
			}

	log.Printf("user %s is added", req.Username)

	render.JSON(w, r, serv.ToResponse(req.Username))

	}
}