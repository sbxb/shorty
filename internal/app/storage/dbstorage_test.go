package storage_test

// import (
// 	"testing"

// 	"github.com/sbxb/shorty/internal/app/storage"
// 	u "github.com/sbxb/shorty/internal/app/url"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// const dsn = "postgres://shorty:shorty@localhost/shortytest"

// func TestDBStorage_AddBatch(t *testing.T) {
// 	batch := []u.BatchURLEntry{
// 		{
// 			CorrelationID: "123",
// 			OriginalURL:   "http://example.com",
// 			ShortURL:      "/5agFZWrIb6Ej21QvYUNBL3",
// 		},
// 		{
// 			CorrelationID: "456",
// 			OriginalURL:   "http://example.org",
// 			ShortURL:      "/6EH6vwAy9dOyyNbopTS6M4",
// 		},
// 	}

// 	store, err := storage.NewDBStorage(dsn)
// 	require.NoError(t, err)
// 	_ = store.Truncate()

// 	err = store.AddBatchURL(batch, "")
// 	require.NoError(t, err)
// }

// func TestDBStorage_Add_then_Get(t *testing.T) {
// 	urls := []string{
// 		"http://example.com",
// 		"http://example.org",
// 	}

// 	store, err := storage.NewDBStorage(dsn)
// 	require.NoError(t, err)
// 	_ = store.Truncate()

// 	for _, url := range urls {
// 		id := u.ShortID(url)

// 		err := store.AddURL(url, id, "")
// 		require.NoError(t, err)

// 		urlReturned, err := store.GetURL(id)
// 		require.NoError(t, err)

// 		assert.Equal(t, urlReturned, url)
// 	}
// }

// func TestDBStorage_Get_Nonexistent(t *testing.T) {
// 	id := "nonexistent_id"

// 	store, err := storage.NewDBStorage(dsn)
// 	require.NoError(t, err)
// 	_ = store.Truncate()

// 	urlReturned, err := store.GetURL(id)
// 	require.NoError(t, err)

// 	assert.Empty(t, urlReturned)
// }

// func TestDBStorage_Add_Record_Twice(t *testing.T) {
// 	url := "http://example.com"

// 	store, err := storage.NewDBStorage(dsn)
// 	require.NoError(t, err)
// 	_ = store.Truncate()

// 	id := u.ShortID(url)
// 	_ = store.AddURL(url, id, "")   // once
// 	err = store.AddURL(url, id, "") // twice

// 	var conflictError *storage.IDConflictError
// 	require.ErrorAs(t, err, &conflictError)
// }
