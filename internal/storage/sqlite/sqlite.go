package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server-2/internal/storage"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		age INTEGER);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db : db}, nil
}

func (s *Storage) Create(username, password string, firstName, lastName *string, age *int) (error) {
	const fn = "storage.sqlite.CreateUser"

	hashedPassword, err := hashPass(password)
	if err != nil {
		return fmt.Errorf("failed to hash pass: %w", err)
	}

	stmt, err := s.db.Prepare("INSERT INTO users(username, password, first_name, last_name, age) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	var fnVal, lnVal interface{} = nil, nil
	var aVal interface{} = nil

	if firstName != nil {
		fnVal = *firstName
	}

	if lastName != nil {
		lnVal = *lastName
	}

	if age != nil {
		aVal = *age
	}


	res, err := stmt.Exec(username, hashedPassword, fnVal, lnVal, aVal)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", fn, storage.ErrUserExists)
		}

		return fmt.Errorf("%s, %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: failed to get last insert id: %w", fn, err)
	}

	log.Printf("User %s with id %v is created", username, id)

	return nil

}

func (s *Storage) GetUser(username string) (*string, *string, *int, error) {
	const fn = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT first_name, last_name, age FROM users WHERE username = ?")
	if err != nil {
		return  nil, nil, nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var firstName, lastName sql.NullString
	var age sql.NullInt64

	err = stmt.QueryRow(username).Scan(&firstName, &lastName, &age)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	var fnIf *string
	if firstName.Valid {
		fnIf = &firstName.String
	}

	var lnIf *string
	if lastName.Valid {
		lnIf = &lastName.String
	}

	var aIf *int
	if age.Valid {
		a := int(age.Int64)
		aIf = &a
	}

	return fnIf, lnIf, aIf, nil
}


func hashPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func (s *Storage) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var hashedPassword string
		err := s.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
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

