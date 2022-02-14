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

	router.Get("/{id}", gzipMW(authMW(urlHandler.GetHandler)))
	router.Post("/", gzipMW(authMW(urlHandler.PostHandler)))

	router.Post("/api/shorten", gzipMW(authMW(urlHandler.JSONPostHandler)))
	router.Post("/api/shorten/batch", gzipMW(authMW(urlHandler.JSONBatchPostHandler)))

	router.Delete("/api/user/urls", gzipMW(authMW(urlHandler.UserDeleteHandler)))

	router.Get("/user/urls", gzipMW(authMW(urlHandler.UserGetHandler)))

	router.Get("/ping", gzipMW(urlHandler.PingGetHandler)) // no auth cookie needed

	return router
}
