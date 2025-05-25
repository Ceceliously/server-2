package main

import (
	"log"
	"os"
	"net/http"
	
	"server-2/internal/config"
	"server-2/internal/storage/sqlite"
	"server-2/internal/http-server/handlers/create"
	"server-2/internal/http-server/handlers/read"

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

	_=storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)



    router.Get("/user", storage.BasicAuth(read.GetUser(storage)))

	router.Post("/user", create.New(storage))


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

