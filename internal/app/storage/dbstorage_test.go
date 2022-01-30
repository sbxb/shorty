package storage_test

// import (
// 	"testing"

// 	"github.com/sbxb/shorty/internal/app/storage"
// 	u "github.com/sbxb/shorty/internal/app/url"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// const dsn = "postgres://shorty:shorty@localhost/shortytest"

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
