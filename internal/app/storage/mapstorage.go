package storage

import (
	"fmt"
	"sync"
)

// MapStorage defines a simple in-memory storage implemented as a wrapper aroung Go map
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

// AddURL saves url and its id
// MapStorage implementation never returns non-nil error
func (st *MapStorage) AddURL(url string, id string) error {
	st.mu.Lock()
	st.data[id] = url
	st.mu.Unlock()

	return nil
}

// GetURL searches for url by its id
// MapStorage implementation never returns non-nil error
func (st *MapStorage) GetURL(id string) (string, error) {
	st.mu.RLock()
	url, ok := st.data[id]
	st.mu.RUnlock()
	if ok {
		return url, nil
	}

	return "", fmt.Errorf("id %s not found", id)
}
