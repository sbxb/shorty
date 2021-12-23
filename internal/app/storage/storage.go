package storage

import (
	"fmt"
	"strconv"
)

type MapStorage struct {
	m map[string]string
}

func NewMapStorage() *MapStorage {
	m := make(map[string]string)
	return &MapStorage{m}
}

func (st *MapStorage) Add(url string) (id string) {
	id = strconv.Itoa(len(st.m) + 1)
	st.m[id] = url
	return id
}

func (st *MapStorage) Get(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	}
	return "", fmt.Errorf("id %s not found", id)
}
