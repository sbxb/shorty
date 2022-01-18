package main

import (
	"log"
	"os"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
)

func main() {
	cfg := config.New()
	if err := cfg.Validate(); err != nil {
		log.Fatalln(err)
	}

	store := storage.NewMapStorage()
	if err := store.Open(cfg.FileStoragePath); err != nil {
		log.Fatalln(err)
	}

	router := api.NewRouter(store, cfg)
	server, _ := api.NewHTTPServer(cfg.ServerAddress, router, store.Close)

	go server.WaitForInterrupt()
	os.Exit(server.Run())
}
