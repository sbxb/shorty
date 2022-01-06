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

func (st *MapStorage) AddURL(url string, id string) error {
	st.mu.Lock()
	st.data[id] = url
	st.mu.Unlock()

	return nil
}

func (st *MapStorage) GetURL(id string) (string, error) {
	st.mu.RLock()
	url, ok := st.data[id]
	st.mu.RUnlock()
	if ok {
		return url, nil
	}

	return "", fmt.Errorf("id %s not found", id)
}
