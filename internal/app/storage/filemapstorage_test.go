package storage_test

import (
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMapStorage_Write_And_Read_File(t *testing.T) {
	tmpFileName := t.TempDir() + "/" + "test.db"

	store, err := storage.NewFileMapStorage(tmpFileName)

	require.NoError(t, err)

	urls := []string{
		"http://example.com",
		"http://example.org",
	}

	for _, url := range urls {
		id := u.ShortID(url)
		_ = store.AddURL(url, id, "") // MapStorage never returns non-nil error
	}

	store.Close()

	store, err = storage.NewFileMapStorage(tmpFileName)

	require.NoError(t, err)

	for _, url := range urls {
		id := u.ShortID(url)
		urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error
		assert.Equal(t, urlReturned, url)
	}

	store.Close()
}
