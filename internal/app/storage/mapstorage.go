package storage

import (
	"sync"
)

// MapStorage defines a simple in-memory storage implemented as a wrapper
// aroung Go map
type MapStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

// MapStorage implements Storage interface
var _ Storage = (*MapStorage)(nil)

func NewMapStorage() *MapStorage {
	d := make(map[string]string)
	return &MapStorage{data: d}
}

// AddURL saves both url and its id
// MapStorage implementation never returns non-nil error
func (st *MapStorage) AddURL(url string, id string) error {
	st.mu.Lock()
	st.data[id] = url
	st.mu.Unlock()

	return nil
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
// MapStorage implementation never returns non-nil error
func (st *MapStorage) GetURL(id string) (string, error) {
	st.mu.RLock()
	url := st.data[id]
	st.mu.RUnlock()

	return url, nil
}
