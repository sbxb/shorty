package storage_test

import (
	"context"
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMapStorage_Write_And_Read_File(t *testing.T) {
	tmpFileName := t.TempDir() + "/" + "test.db"

	store, err := storage.NewFileMapStorage(tmpFileName)

	require.NoError(t, err)

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

	for _, ue := range entries {
		err = store.AddURL(context.Background(), ue, "")
		require.NoError(t, err)
	}

	store.Close()

	// Reading all the items written previously
	store, err = storage.NewFileMapStorage(tmpFileName)

	require.NoError(t, err)

	for _, ue := range entries {
		urlReturned, _ := store.GetURL(context.Background(), ue.ShortURL) // MapStorage.GetURL() never returns non-nil error
		assert.Equal(t, urlReturned, ue.OriginalURL)
	}

	store.Close()
}
