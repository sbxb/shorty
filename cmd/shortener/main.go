package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sbxb/shorty/internal/app/api"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	var db *sql.DB
	if cfg.DatabaseDSN != "" {
		db, err = sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			log.Println("Server failed to open DB: " + err.Error())
			db = nil
		} else {
			defer db.Close()
		}
		handlers.Database = db
	}

	store, err := storage.NewFileMapStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	router := api.NewRouter(store, cfg)
	server, err := api.NewHTTPServer(cfg.ServerAddress, router)
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
}
