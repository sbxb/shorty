package main

import (
	"net/http"

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

	http.ListenAndServe(cfg.FullServerName(), router)
}
