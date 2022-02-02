package storage_test

import (
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_Add_then_Get(t *testing.T) {
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://local.test",
	}

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error

	for _, url := range urls {
		id := u.ShortID(url)
		err := store.AddURL(url, id, "")
		require.NoError(t, err)
		urlReturned, _ := store.GetURL(id) // MapStorage.GetURL() never returns non-nil error

		assert.Equal(t, urlReturned, url)
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	id := "nonexistent_id"

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error

	urlReturned, _ := store.GetURL(id) // MapStorage.GetURL() never returns non-nil error

	assert.Empty(t, urlReturned)
}

func TestMemoryStore_Add_Record_Twice(t *testing.T) {
	url := "http://example.com"

	store, _ := storage.NewMapStorage() // NewMapStorage() never returns non-nil error
	id := u.ShortID(url)
	_ = store.AddURL(url, id, "")    // once
	err := store.AddURL(url, id, "") // twice

	var conflictError *storage.IDConflictError
	require.ErrorAs(t, err, &conflictError)
}
