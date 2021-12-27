package main

import (
	"net/http"

	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

const serverName = "localhost:8080"

func main() {
	r := chi.NewRouter()

	store := storage.NewMapStorage()

	// r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
	// 	rw.Write([]byte("chi"))
	// })

	r.Get("/{id}", handlers.GetHandler(store, serverName))
	r.Post("/", handlers.PostHandler(store, serverName))

	http.ListenAndServe(serverName, r)
}
