package api

import (
	"net/http"

	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

// NewRouter creates chi router and handlers container, register handlers and
// pass dependencies to handlers
func NewRouter(store storage.Storage, cfg config.Config) http.Handler {
	router := chi.NewRouter()

	urlHandler := handlers.NewURLHandler(store, cfg)

	router.Get("/{id}", urlHandler.GetHandler)
	router.Post("/", urlHandler.PostHandler)

	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	return router
}
