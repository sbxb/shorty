package storage

import (
	"context"

	"github.com/sbxb/shorty/internal/app/url"
)

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	AddURL(ctx context.Context, ue url.URLEntry, userID string) error
	AddBatchURL(ctx context.Context, batch []url.BatchURLEntry, userID string) error
	// TODO return URLEntry instead of a string
	GetURL(id string) (string, error)
	GetUserURLs(userID string) ([]url.URLEntry, error)
	Close() error
}
