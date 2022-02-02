package storage

import (
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
func (st *MapStorage) AddURL(url string, id string, uid string) error {
	st.Lock()
	defer st.Unlock()
	if uid == "" {
		uid = "NULL"
	}
	if _, ok := st.data[id]; ok {
		log.Println(">>> MapStorage: Repeated id found: ", id)
		return NewIDConflictError(id)
	}
	st.data[id] = uid + ";" + url

	return nil
}

func (st *MapStorage) AddBatchURL(batch []url.BatchURLEntry, uid string) error {
	st.Lock()
	defer st.Unlock()
	if uid == "" {
		uid = "NULL"
	}
	for _, e := range batch {
		st.data[e.ShortURL] = uid + ";" + e.OriginalURL
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

func (st *MapStorage) GetUserURLs(uid string) ([]url.UserURL, error) {
	res := []url.UserURL{}
	for id, str := range st.data {
		parts := strings.SplitN(str, ";", 2)
		if parts[0] != uid {
			continue
		}
		entry := url.UserURL{
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
