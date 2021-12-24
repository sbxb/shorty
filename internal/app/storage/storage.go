package storage

import (
	"fmt"
	"strconv"
)

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	AddURL(url string) string
	GetURL(id string) (string, error)
}

// MapStorage defines a simple in-memory storage implemented as a wrapper aroung Go map
// FIXME not thread-safe at the moment
type MapStorage struct {
	m map[string]string
}

// MapStorage implements Storage interface
var _ Storage = (*MapStorage)(nil)

func NewMapStorage() *MapStorage {
	m := make(map[string]string)
	return &MapStorage{m}
}

func (st *MapStorage) AddURL(url string) (id string) {
	id = strconv.Itoa(len(st.m) + 1)
	st.m[id] = url
	return id
}

func (st *MapStorage) GetURL(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	}
	return "", fmt.Errorf("id %s not found", id)
}
