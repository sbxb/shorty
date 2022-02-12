package storage

import (
	"context"
	"strings"
	"sync"

	"github.com/sbxb/shorty/internal/app/logger"
	"github.com/sbxb/shorty/internal/app/url"
)

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

	if _, ok := st.data[ue.ShortURL]; ok {
		logger.Info("MapStorage: Repeated id found: ", ue.ShortURL)
		return NewIDConflictError(ue.ShortURL)
	}
	st.data[ue.ShortURL] = userID + "|false|" + ue.OriginalURL
	logger.Debugf("AddURL [%s] :: [%s]", ue.ShortURL, st.data[ue.ShortURL])

	return nil
}

func (st *MapStorage) AddBatchURL(ctx context.Context, batch []url.BatchURLEntry, userID string) error {
	st.Lock()
	defer st.Unlock()

	for _, ue := range batch {
		st.data[ue.ShortURL] = userID + "|false|" + ue.OriginalURL
	}

	return nil
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
// MapStorage implementation never returns non-nil error (except for records
// marked as deleted)
func (st *MapStorage) GetURL(ctx context.Context, id string) (string, error) {
	st.RLock()
	defer st.RUnlock()
	res := st.data[id]
	if res == "" {
		return res, nil
	}
	parts := strings.SplitN(res, "|", 3)
	if parts[1] == "true" {
		logger.Info("MapStorage: Deleted id found: ", id)
		return "", NewURLDeletedError(id)
	}

	return parts[2], nil
}

func (st *MapStorage) GetUserURLs(ctx context.Context, userID string) ([]url.URLEntry, error) {
	res := []url.URLEntry{}
	for id, str := range st.data {
		parts := strings.SplitN(str, "|", 3)
		if parts[0] != userID {
			continue
		}
		entry := url.URLEntry{
			ShortURL:    id,
			OriginalURL: parts[2],
		}
		res = append(res, entry)
	}
	return res, nil
}

func (st *MapStorage) DeleteBatch(ctx context.Context, ids []string, userID string) error {
	logger.Debugf("MapStorage : DeleteBatch: Got ids %v", ids)
	st.Lock()
	defer st.Unlock()

	idHash := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idHash[id] = struct{}{}
	}

	for id, str := range st.data {
		if _, ok := idHash[id]; !ok {
			continue
		}
		parts := strings.SplitN(str, "|", 3)
		if parts[0] != userID || parts[1] != "false" {
			logger.Debugf("MapStorage : DeleteBatch: skip id %s", id)
			continue
		}
		st.data[id] = parts[0] + "|true|" + parts[2]
		logger.Debugf("MapStorage : DeleteBatch: id %s marked deleted", id)
	}
	return nil
}

func (st *MapStorage) Close() error {
	return nil
}
