package read

import (
	"log"
	"net/http"
	"errors"
	"server-2/internal/storage"
	res "server-2/internal/lib/user/response"

	"github.com/go-chi/render"

)



type Response struct {
	res.Response
	res.User
}




func GetUser(s storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.read"

		username, _, _ := r.BasicAuth()


	firstName, lastName, age, err := s.GetUser(username)
	if errors.Is(err, storage.ErrUserNotFound) {
		log.Printf("failed to find user : %w", err)
		render.JSON(w, r, res.Error("failed to find user"))
		return
	}
	if err != nil {
			log.Printf("%s: failed to find user", fn)
			render.JSON(w, r, res.Error("failed to find user"))
			return
		}

		resp := Response {
		Response: res.OK(),
		User: res.User {
			Username: username,
		},
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

	log.Printf("user %s is found", username)

	render.JSON(w, r, resp)
	}
}