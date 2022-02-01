package storage

import (
	"github.com/sbxb/shorty/internal/app/url"
)

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	AddURL(url string, id string, uid string) error
	AddBatchURL(batch []url.BatchURLEntry, uid string) error
	GetURL(id string) (string, error)
	GetUserURLs(uid string) ([]url.UserURL, error)
	Close() error
}
