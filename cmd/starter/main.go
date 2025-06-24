package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"server-2/internal/config"
	"server-2/internal/middleware/auth"
	"server-2/internal/service/user_service"
	"server-2/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, stop := context.WithTimeout(context.Background(), 30*time.Second)
		defer stop()


	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Println("failed to init storage", err)
		os.Exit(1)
	}


	userService := user_service.NewUserService(storage)
	basicAuth := auth.NewBasicAuth(storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)



	router.Post("/user", userService.HandlersV1.CreateUserHandler)
	router.Get("/user", basicAuth.BasicAuth(userService.HandlersV1.GetUserHandler))

	srv := &http.Server{
        Addr:         cfg.Address,
        Handler:      router,
        IdleTimeout:  cfg.HTTPServer.IdleTimeout,
        ReadTimeout:  cfg.HTTPServer.Timeout,
        WriteTimeout: cfg.HTTPServer.Timeout,
    }

	go func () {
		log.Println("starting server on ", cfg.Address)
		err = srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	} ()

	log.Printf("listening on %s", cfg.Address)

	sig := <-signalChan 
		log.Printf("recieve shutdown signal: %v.", sig)

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown error")
	}

	log.Println("shutting down server gracefully")
	
	

}

