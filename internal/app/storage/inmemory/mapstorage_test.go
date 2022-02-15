package inmemory_test

import (
	"context"
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/storage/inmemory"
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

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	for _, ue := range entries {
		err := store.AddURL(context.Background(), ue, "")
		require.NoError(t, err)
		urlReturned, err := store.GetURL(context.Background(), ue.ShortURL)
		require.NoError(t, err)

		assert.Equal(t, urlReturned, ue.OriginalURL)
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	id := "nonexistent_id"

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	urlReturned, err := store.GetURL(context.Background(), id)
	require.NoError(t, err)

	assert.Empty(t, urlReturned)
}

func TestMemoryStore_Add_Record_Twice(t *testing.T) {
	ue := url.URLEntry{
		ShortURL:    "5agFZWrIb6Ej21QvYUNBL3",
		OriginalURL: "http://example.com",
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	ctx := context.Background()
	_ = store.AddURL(ctx, ue, "")    // once
	err := store.AddURL(ctx, ue, "") // twice

	var conflictError *storage.IDConflictError
	require.ErrorAs(t, err, &conflictError)
}

func TestMemoryStore_Batch_Add_Delete(t *testing.T) {
	entries := []url.BatchURLEntry{
		{
			ShortURL:    "5agFZWrIb6Ej21QvYUNBL3",
			OriginalURL: "http://example.com",
		},
		{
			ShortURL:    "6EH6vwAy9dOyyNbopTS6M4",
			OriginalURL: "http://example.org",
		},
	}

	store, _ := inmemory.NewMapStorage() // NewMapStorage() never returns non-nil error

	err := store.AddBatchURL(context.Background(), entries, "")
	require.NoError(t, err)

	for _, ue := range entries {
		urlReturned, _ := store.GetURL(context.Background(), ue.ShortURL)
		require.NoError(t, err)

		assert.Equal(t, urlReturned, ue.OriginalURL)
	}

	ids := make([]string, len(entries))
	for _, ue := range entries {
		ids = append(ids, ue.ShortURL)
	}

	err = store.DeleteBatch(context.Background(), ids, "")
	require.NoError(t, err)

	for _, ue := range entries {
		_, err = store.GetURL(context.Background(), ue.ShortURL)
		var deletedError *storage.URLDeletedError
		require.ErrorAs(t, err, &deletedError)
	}
}
