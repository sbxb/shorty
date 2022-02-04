package storage_test

import (
	"context"
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_Add_then_Get(t *testing.T) {
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

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error

	for _, ue := range entries {
		err := store.AddURL(context.Background(), ue, "")
		require.NoError(t, err)
		urlReturned, _ := store.GetURL(ue.ShortURL) // MapStorage.GetURL() never returns non-nil error

		assert.Equal(t, urlReturned, ue.OriginalURL)
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	id := "nonexistent_id"

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error

	urlReturned, _ := store.GetURL(id) // MapStorage.GetURL() never returns non-nil error

	assert.Empty(t, urlReturned)
}

func TestMemoryStore_Add_Record_Twice(t *testing.T) {
	ue := url.URLEntry{
		ShortURL:    "5agFZWrIb6Ej21QvYUNBL3",
		OriginalURL: "http://example.com",
	}

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error

	ctx := context.Background()
	_ = store.AddURL(ctx, ue, "")    // once
	err := store.AddURL(ctx, ue, "") // twice

	var conflictError *storage.IDConflictError
	require.ErrorAs(t, err, &conflictError)
}
