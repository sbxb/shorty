package main

import (
	"os"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
)

func main() {
	cfg := config.DefaultConfig
	store := storage.NewMapStorage()
	router := api.NewRouter(store, cfg.FullServerName())
	server, _ := api.NewHTTPServer(cfg.FullServerName(), router)

	go server.WaitForInterrupt()
	os.Exit(server.Run())
}
