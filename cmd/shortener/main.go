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

	// TODO:
	// mux := router.NewRouter()
	// http.ListenAndServe(serverName, r)
	r := chi.NewRouter()

	store := storage.NewMapStorage()

	r.Get("/{id}", handlers.GetHandler(store, cfg.FullServerName()))
	r.Post("/", handlers.PostHandler(store, cfg.FullServerName()))

	http.ListenAndServe(cfg.FullServerName(), r)
}
