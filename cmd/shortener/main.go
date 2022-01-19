package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
)

func main() {
	var wg sync.WaitGroup

	cfg := config.New()
	if err := cfg.Validate(); err != nil {
		log.Fatalln(err)
	}

	store := storage.NewMapStorage()
	if err := store.Open(cfg.FileStoragePath); err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	router := api.NewRouter(store, cfg)
	server, _ := api.NewHTTPServer(cfg.ServerAddress, router)

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGTERM, syscall.SIGINT,
	)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start(ctx)
	}()

	wg.Wait()
}
