package storage_test

import (
	"testing"

	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"
)

func TestMemoryStore_Add_then_Get(t *testing.T) {
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://local.test",
	}

	store := storage.NewMapStorage()

	for _, url := range urls {
		id := u.ShortId(url)
		_ = store.AddURL(url, id)          // MapStorage never returns non-nil error
		urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error
		if urlReturned != url {
			t.Errorf("got [%s], want [%s]", urlReturned, url)
		}
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	store := storage.NewMapStorage()
	id := "nonexistent_id"
	urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error
	if urlReturned != "" {
		t.Errorf("got [%s] which is not supposed to exist in storage", urlReturned)
	}
}
