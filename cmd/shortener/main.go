package main

import (
	"os"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
)

func main() {
	cfg := config.New()
	store := storage.NewMapStorage()
	router := api.NewRouter(store, cfg)
	server, _ := api.NewHTTPServer(cfg.ServerAddress, router)

	go server.WaitForInterrupt()
	os.Exit(server.Run())
}
