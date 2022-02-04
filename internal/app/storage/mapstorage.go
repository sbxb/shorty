package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sbxb/shorty/internal/app/url"
)

type IDConflictError struct {
	ID string
}

func (ice *IDConflictError) Error() string {
	return fmt.Sprintf("Storage already has a record with id %s", ice.ID)
}

func NewIDConflictError(id string) error {
	return &IDConflictError{id}
}

// MapStorage defines a simple in-memory storage implemented as a wrapper
// around Go map
type MapStorage struct {
	sync.RWMutex

	data map[string]string
}

// MapStorage implements Storage interface
var _ Storage = (*MapStorage)(nil)

func NewMapStorage() (*MapStorage, error) {
	d := make(map[string]string)
	return &MapStorage{data: d}, nil
}

// AddURL saves both url and its id
func (st *MapStorage) AddURL(ctx context.Context, ue url.URLEntry, userID string) error {
	st.Lock()
	defer st.Unlock()
	if userID == "" {
		userID = "NULL"
	}
	if _, ok := st.data[ue.ShortURL]; ok {
		log.Println(">>> MapStorage: Repeated id found: ", ue.ShortURL)
		return NewIDConflictError(ue.ShortURL)
	}
	st.data[ue.ShortURL] = userID + ";" + ue.OriginalURL

	return nil
}

func (st *MapStorage) AddBatchURL(ctx context.Context, batch []url.BatchURLEntry, userID string) error {
	st.Lock()
	defer st.Unlock()
	// TODO check if empty string could be used in split at GetURL()
	if userID == "" {
		userID = "NULL"
	}
	for _, ue := range batch {
		st.data[ue.ShortURL] = userID + ";" + ue.OriginalURL
	}

	return nil
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
// MapStorage implementation never returns non-nil error
func (st *MapStorage) GetURL(id string) (string, error) {
	st.RLock()
	defer st.RUnlock()
	res := st.data[id]
	if res == "" {
		return res, nil
	}
	parts := strings.SplitN(res, ";", 2)

	return parts[1], nil
}

func (st *MapStorage) GetUserURLs(userID string) ([]url.URLEntry, error) {
	res := []url.URLEntry{}
	for id, str := range st.data {
		parts := strings.SplitN(str, ";", 2)
		if parts[0] != userID {
			continue
		}
		entry := url.URLEntry{
			ShortURL:    id,
			OriginalURL: parts[1],
		}
		res = append(res, entry)
	}
	return res, nil
}

func (st *MapStorage) Close() error {
	return nil
}
