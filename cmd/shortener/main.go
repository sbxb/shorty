package main

import (
	"net/http"

	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"
)

const serverName = "localhost:8080"

func main() {
	store := storage.NewMapStorage()
	http.Handle("/", handlers.DefaultHandler(store, serverName))

	http.ListenAndServe(serverName, nil)
}
