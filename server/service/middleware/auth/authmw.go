package auth


import (
	"context"
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"server-2/internal/storage"
)


type BasicAuth struct {
	storage storage.UserStorage
}

func NewBasicAuth(s storage.UserStorage) *BasicAuth {
	return &BasicAuth{storage: s}
}

func (a *BasicAuth) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := a.storage.GetPasswordHash(username)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
				ctx := context.WithValue(r.Context(), "username", username)
				next.ServeHTTP(w, r.WithContext(ctx))
		
	}
}