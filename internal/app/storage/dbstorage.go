package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DBStorage defines a simple in-memory storage implemented as a wrapper
// aroung Go map
type DBStorage struct {
	db *sql.DB
}

// DBStorage implements Storage interface
//var _ Storage = (*DBStorage)(nil)

const defaultTimeout = 5 * time.Second

func NewDBStorage(dsn string) (*DBStorage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DBStorage: empty dsn")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("DBStorage: Open: %v", err)
	}

	// Let's ping the database before returning DBStorage instance
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Ping: %v", err)
	}

	// Let's create all the necessary tables in the database
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Create Tables: %v", err)
	}

	return &DBStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	URLsTableName := "urls"
	URLsTableQuery := `CREATE TABLE IF NOT EXISTS ` + URLsTableName + ` (
		id INT primary key GENERATED ALWAYS AS IDENTITY,
		url_id VARCHAR(512) NOT NULL,
		user_id VARCHAR(512) NOT NULL,
		original_url TEXT NOT NULL
	)`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if _, err := db.ExecContext(ctx, URLsTableQuery); err != nil {
		return err
	}

	return nil
}

// Ping pings the database
func (st *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if err := st.db.PingContext(ctx); err != nil {
		return fmt.Errorf("DBStorage: %v", err)
	}
	return nil
}

func (st *DBStorage) Close() error {
	if st.db == nil {
		return nil
	}

	if err := st.db.Close(); err != nil {
		return err
	}

	st.db = nil

	return nil
}
