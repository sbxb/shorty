package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/logger"
	"github.com/sbxb/shorty/internal/app/storage"
)

func main() {
	var wg sync.WaitGroup

	logger.SetLevel("WARNING")

	cfg, err := config.New()
	if err != nil {
		logger.Fatalln(err)
	}

	var store storage.Storage

	if cfg.DatabaseDSN != "" {
		store, err = storage.NewDBStorage(cfg.DatabaseDSN)
	} else {
		store, err = storage.NewFileMapStorage(cfg.FileStoragePath)
	}
	if err != nil {
		logger.Fatalln(err)
	}
	defer store.Close()

	router := api.NewRouter(store, cfg)
	server, err := api.NewHTTPServer(cfg.ServerAddress, router)
	if err != nil {
		logger.Fatalln(err)
	}
	defer server.Close()

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
	if err := store.Close(); err != nil {
		logger.Error(err)
	}
}
