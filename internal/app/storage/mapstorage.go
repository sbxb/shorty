package storage

import (
	"fmt"
	"strconv"
)

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
	// TODO There should be more sophisticated algorithm for calculating URL's short id
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
