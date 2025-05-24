package read

import (
	"log"
	"net/http"
	"errors"
	"server-2/internal/storage"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Resp struct {
    Status string `json:"status"`
    Error  string `json:"error,omitempty"`
}

type Response struct {
	Resp
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Age int `json:"age"`
}

const (
	StatusOk = "OK"
	StatusError = "Error"
)

func OK() Resp {
	return Resp{
		Status: StatusOk,
	} 
}

func Error(msg string) Resp {
	return Resp{
		Status: StatusError,
		Error: msg,
	}
}


func GetUser(s storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.read"

		username, password, ok := r.BasicAuth()
      if !ok {
          w.WriteHeader(http.StatusBadRequest)
		  render.JSON(w,r,Error("auth required"))
          return
      }

		req := Request {
			Username: username,
			Password: password,
		}

	if err := validator.New().Struct(req); err != nil {
		log.Printf("invalid request : %w", err)
		render.JSON(w, r, Error("invalid request"))
		return
	}

	firstName, lastName, age, err := s.GetUser(username, password)
	if errors.Is(err, storage.ErrUserNotFound) {
		log.Printf("failed to find user : %w", err)
		render.JSON(w, r, Error("failed to find user"))
		return
	}
	if err != nil {
			log.Printf("%s: failed to find user", fn)
			render.JSON(w, r, Error("failed to find user"))
			return
		}

		resp := Response {
		Resp: OK(),
		Username: username,
	}

	if firstName != nil {
		resp.FirstName = *firstName
	}

	if lastName != nil {
		resp.LastName = *lastName
	}

	if age != nil {
		resp.Age = *age
	}

	log.Printf("user %s is found", req.Username)

	render.JSON(w, r, resp)
	}
}