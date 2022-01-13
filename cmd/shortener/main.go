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
	store := storage.NewMapStorage()
	if err := store.BindFile(cfg.FileStoragePath); err != nil {
		log.Fatalln(err)
	}
	router := api.NewRouter(store, cfg)
	server, _ := api.NewHTTPServer(cfg.ServerAddress, router, store.Close)

	go server.WaitForInterrupt()
	os.Exit(server.Run())
}
