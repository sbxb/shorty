package main

import (
	"os"

	"github.com/sbxb/shorty/internal/app/api"
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

	server, _ := api.NewHTTPServer(cfg.FullServerName(), router)
	go server.WaitForInterrupt()
	os.Exit(server.Run())
}
