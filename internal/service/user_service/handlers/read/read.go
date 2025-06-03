package read

import (
	"log"
	"net/http"
	"errors"
	"server-2/internal/storage"
	res "server-2/internal/models/response"
	serv "server-2/internal/models/user/user_get"
	"server-2/internal/service/user_service/usecase/user"

	"github.com/go-chi/render"

)



func GetUserHandler(uc *user.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.read"

	username, ok := r.Context().Value("username").(string)
		if !ok {
			log.Printf("%s: username not found in context", fn)
			render.JSON(w, r, res.Error("internal server error"))
			return
		}

	ObtainedUser, err := uc.GetUser(username)
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

	log.Printf("user %s is found", username)

	render.JSON(w, r, serv.ToResponse(ObtainedUser))
	}
}