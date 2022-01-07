package api

import (
	"net/http"

	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func NewRouter(store storage.Storage, serverName string) http.Handler {
	router := chi.NewRouter()

	urlHandler := handlers.URLHandler{
		Store:      store,
		ServerName: serverName,
	}

	router.Get("/{id}", urlHandler.GetHandler)
	router.Post("/", urlHandler.PostHandler)

	return router
}
