package storage_test

import (
	"context"
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dsn = "postgres://shorty:shorty@localhost/shortytest"

func TestDBStorage_AddBatch(t *testing.T) {
	batch := []url.BatchURLEntry{
		{
			CorrelationID: "123",
			OriginalURL:   "http://example.com",
			ShortURL:      "/5agFZWrIb6Ej21QvYUNBL3",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "http://example.org",
			ShortURL:      "/6EH6vwAy9dOyyNbopTS6M4",
		},
	}

	store, err := storage.NewDBStorage(dsn)
	require.NoError(t, err)
	_ = store.Truncate()

	err = store.AddBatchURL(context.Background(), batch, "")
	require.NoError(t, err)
}

func TestDBStorage_Add_then_Get(t *testing.T) {
	entries := []url.URLEntry{
		{
			ShortURL:    "5agFZWrIb6Ej21QvYUNBL3",
			OriginalURL: "http://example.com",
		},
		{
			ShortURL:    "6EH6vwAy9dOyyNbopTS6M4",
			OriginalURL: "http://example.org",
		},
	}

	store, err := storage.NewDBStorage(dsn)
	require.NoError(t, err)
	_ = store.Truncate()

	for _, ue := range entries {
		err := store.AddURL(context.Background(), ue, "")
		require.NoError(t, err)

		urlReturned, err := store.GetURL(ue.ShortURL)
		require.NoError(t, err)

		assert.Equal(t, urlReturned, ue.OriginalURL)
	}
}

func TestDBStorage_Get_Nonexistent(t *testing.T) {
	id := "nonexistent_id"

	store, err := storage.NewDBStorage(dsn)
	require.NoError(t, err)
	_ = store.Truncate()

	urlReturned, err := store.GetURL(id)
	require.NoError(t, err)

	assert.Empty(t, urlReturned)
}

func TestDBStorage_Add_Record_Twice(t *testing.T) {
	ue := url.URLEntry{
		ShortURL:    "5agFZWrIb6Ej21QvYUNBL3",
		OriginalURL: "http://example.com",
	}

	store, err := storage.NewDBStorage(dsn)
	require.NoError(t, err)
	_ = store.Truncate()

	ctx := context.Background()

	_ = store.AddURL(ctx, ue, "")   // once
	err = store.AddURL(ctx, ue, "") // twice

	var conflictError *storage.IDConflictError
	require.ErrorAs(t, err, &conflictError)
}
