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

	router.Use(gzipMW)
	router.Use(authMW)

	router.Get("/{id}", urlHandler.GetHandler)
	router.Post("/", urlHandler.PostHandler)

	router.With(jsonEncMW).Post("/api/shorten", urlHandler.JSONPostHandler)
	router.With(jsonEncMW).Post("/api/shorten/batch", urlHandler.JSONBatchPostHandler)

	router.With(jsonEncMW).Delete("/api/user/urls", urlHandler.UserDeleteHandler)

	router.Get("/api/user/urls", urlHandler.UserGetHandler)

	router.Get("/ping", urlHandler.PingGetHandler)

	return router
}
