package main

import (
	"log"
	"net/http"
	"os"

	"server-2/internal/config"
	"server-2/internal/http-server/handlers/create"
	"server-2/internal/http-server/handlers/read"
	"server-2/internal/storage/sqlite"
	"server-2/server/service/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

		storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Println("failed to init storage", err)
		os.Exit(1)
	}

	userUC := user.NewUserUseCase(storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)



    router.Get("/user", storage.BasicAuth(read.GetUserHandler(userUC)))

	router.Post("/user", create.CreateUserHandler(userUC))


	log.Println("starting server on ", cfg.Address)

	srv := &http.Server{
        Addr:         cfg.Address,
        Handler:      router,
        IdleTimeout:  cfg.HTTPServer.IdleTimeout,
        ReadTimeout:  cfg.HTTPServer.Timeout,
        WriteTimeout: cfg.HTTPServer.Timeout,
    }

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	
}

