package storage

import (
	"sync"
)

// MapStorage defines a simple in-memory storage implemented as a wrapper
// aroung Go map
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
// MapStorage implementation never returns non-nil error
func (st *MapStorage) AddURL(url string, id string) error {
	st.Lock()
	defer st.Unlock()

	st.data[id] = url

	return nil
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
// MapStorage implementation never returns non-nil error
func (st *MapStorage) GetURL(id string) (string, error) {
	st.RLock()
	defer st.RUnlock()
	url := st.data[id]

	return url, nil
}

func (st *MapStorage) Close() error {
	return nil
}
