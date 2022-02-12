package storage

import (
	"context"

	"github.com/sbxb/shorty/internal/app/url"
)

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	AddURL(ctx context.Context, ue url.URLEntry, userID string) error
	AddBatchURL(ctx context.Context, batch []url.BatchURLEntry, userID string) error
	GetURL(ctx context.Context, id string) (string, error)
	GetUserURLs(ctx context.Context, userID string) ([]url.URLEntry, error)
	DeleteBatch(ctx context.Context, ids []string, userID string) error
	Close() error
}
