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

	store := storage.NewMapStorage()

	for _, url := range urls {
		id := u.ShortID(url)
		_ = store.AddURL(url, id)          // MapStorage never returns non-nil error
		urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error

		assert.Equal(t, urlReturned, url)
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	store := storage.NewMapStorage()

	id := "nonexistent_id"
	urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error

	assert.Empty(t, urlReturned)
}

func TestMemoryStore_Save_To_File_And_Read_Again(t *testing.T) {
	tmpFileName := t.TempDir() + "/" + "test.db"

	store := storage.NewMapStorage()
	err := store.Open(tmpFileName)

	require.NoError(t, err)

	urls := []string{
		"http://example.com",
		"http://example.org",
	}

	for _, url := range urls {
		id := u.ShortID(url)
		_ = store.AddURL(url, id) // MapStorage never returns non-nil error
	}

	store.Close()

	store = storage.NewMapStorage()
	err = store.Open(tmpFileName)

	require.NoError(t, err)

	for _, url := range urls {
		id := u.ShortID(url)
		urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error
		assert.Equal(t, urlReturned, url)
	}

	store.Close()
}
