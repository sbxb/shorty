package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.DefaultConfig

	store := storage.NewMapStorage()

	router := chi.NewRouter()
	router.Get("/{id}", handlers.GetHandler(store, cfg.FullServerName()))
	router.Post("/", handlers.PostHandler(store, cfg.FullServerName()))

	// Set more reasonable timeouts than the default ones
	srv := &http.Server{
		Addr:         cfg.FullServerName(),
		Handler:      router,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  36 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
