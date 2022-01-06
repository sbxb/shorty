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
		err := store.AddURL(url, id)
		if err != nil {
			t.Errorf("AddURL() failed")
		}
		urlReturned, err := store.GetURL(id)
		if err != nil {
			t.Errorf("want [%s] but failed to get one", url)
		} else if urlReturned != url {
			t.Errorf("got [%s], want [%s]", urlReturned, url)
		}
	}
}

func TestMemoryStore_Get_Nonexistent(t *testing.T) {
	store := storage.NewMapStorage()
	id := "nonexistent_id"
	urlReturned, err := store.GetURL(id)
	if err == nil {
		t.Errorf("got [%s] which is not supposed to exist in storage", urlReturned)
	}
}
