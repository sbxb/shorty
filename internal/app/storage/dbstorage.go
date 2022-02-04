package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sbxb/shorty/internal/app/url"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// DBStorage defines a database storage implemented as a wrapper
// around any database/sql implementation
type DBStorage struct {
	db       *sql.DB
	urlTable string
}

// DBStorage implements Storage interface
var _ Storage = (*DBStorage)(nil)

const defaultTimeout = 5 * time.Second

func NewDBStorage(dsn string) (*DBStorage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DBStorage: empty dsn")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("DBStorage: Open: %v", err)
	}

	// ping the database before returning DBStorage instance
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Ping: %v", err)
	}

	// create all the necessary tables in the database
	urlTable := "urls"
	if err := createTables(db, urlTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Create Tables: %v", err)
	}

	return &DBStorage{db: db, urlTable: urlTable}, nil
}

func createTables(db *sql.DB, urlTable string) error {
	URLsTableQuery := `CREATE TABLE IF NOT EXISTS ` + urlTable + ` (
		id INT primary key GENERATED ALWAYS AS IDENTITY,
		url_id VARCHAR(512) NOT NULL,
		user_id VARCHAR(512) NOT NULL,
		original_url TEXT NOT NULL,
		UNIQUE (url_id)
	)`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if _, err := db.ExecContext(ctx, URLsTableQuery); err != nil {
		return err
	}

	return nil
}

func (st *DBStorage) Truncate() error {
	URLsTableQuery := `TRUNCATE ` + st.urlTable + ` RESTART IDENTITY`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if _, err := st.db.ExecContext(ctx, URLsTableQuery); err != nil {
		return err
	}

	return nil
}

// AddURL saves both url and its id
func (st *DBStorage) AddURL(ctx context.Context, ue url.URLEntry, userID string) error {
	AddURLQuery := `INSERT INTO ` + st.urlTable + `(url_id, user_id, original_url) 
		VALUES($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result, err := st.db.ExecContext(ctx, AddURLQuery, ue.ShortURL, userID, ue.OriginalURL)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return NewIDConflictError(ue.ShortURL)
		}
		return fmt.Errorf("DBStorage: AddURL: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DBStorage: AddURL: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("DBStorage: AddURL: expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (st *DBStorage) AddBatchURL(ctx context.Context, batch []url.BatchURLEntry, userID string) error {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("DBStorage: AddBatchURL: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO ` + st.urlTable + `(url_id, user_id, original_url)
		VALUES($1, $2, $3) ON CONFLICT(url_id) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("DBStorage: AddBatchURL: %v", err)
	}
	defer stmt.Close()

	for _, e := range batch {
		if _, err = stmt.Exec(e.ShortURL, userID, e.OriginalURL); err != nil {
			return fmt.Errorf("DBStorage: AddBatchURL: %v", err)
		}
	}

	return tx.Commit()
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
func (st *DBStorage) GetURL(id string) (string, error) {
	var url string
	GetURLQuery := `SELECT original_url FROM ` + st.urlTable + ` WHERE 
		url_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := st.db.QueryRowContext(ctx, GetURLQuery, id).Scan(&url)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", fmt.Errorf("DBStorage: GetURL: %v", err)
	default:
		return url, nil
	}
}

func (st *DBStorage) GetUserURLs(userID string) ([]url.URLEntry, error) {
	res := []url.URLEntry{}

	GetUserURLsQuery := `SELECT url_id, original_url FROM ` + st.urlTable + ` WHERE 
		user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	rows, err := st.db.QueryContext(ctx, GetUserURLsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DBStorage: GetUserURLs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u url.URLEntry
		err = rows.Scan(&u.ShortURL, &u.OriginalURL)
		if err != nil {
			return nil, fmt.Errorf("DBStorage: GetUserURLs: %v", err)
		}
		res = append(res, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("DBStorage: GetUserURLs: %v", err)
	}

	return res, nil
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
