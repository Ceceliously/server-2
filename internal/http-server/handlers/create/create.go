package create

import (
	"log"
	"net/http"
	"errors"
	"server-2/internal/storage"
	res "server-2/internal/lib/user/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FirstName *string `json:"first_name,omitempty"`
	LastName *string `json:"last_name,omitempty"`
	Age *int `json:"age,omitempty"`
}



type Response struct {
	res.Response
	Username string `json:"username,omitempty"`
}




func New(s storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.create.new"

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Printf("%s: failed to decode request body", fn)
			render.JSON(w, r, res.Error("failed to decode request"))
			return
		}

	if err := validator.New().Struct(req); err != nil {
		log.Printf("invalid request : %w", err)
		render.JSON(w, r, res.Error("invalid request"))
		return
	}

	err = s.Create(req.Username, req.Password, req.FirstName, req.LastName, req.Age)
	if errors.Is(err, storage.ErrUserExists) {
		log.Printf("failed to safe user : %w", err)
		render.JSON(w, r, res.Error("failed to safe user"))
		return
	}
	if err != nil {
			log.Printf("%s: failed to add user", fn)
			render.JSON(w, r, res.Error("failed to add user"))
			return
		}

	log.Printf("user %s is added", req.Username)

	render.JSON(w, r, Response {
		Response: res.OK(),
		Username: req.Username,
	})
	}
}