package storage_test

import (
	"math/rand"
	"testing"
	"time"

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
		id := u.ShortID(url)
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

func TestMemoryStore_Save_To_File_And_Read_Again(t *testing.T) {
	tmpDirName := t.TempDir()
	tmpFileName := tmpDirName + "/" + getRandomFileName(t)

	store := storage.NewMapStorage()
	if err := store.BindFile(tmpFileName); err != nil {
		t.Fatal(err)
	}

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
	if err := store.BindFile(tmpFileName); err != nil {
		t.Fatal(err)
	}

	for _, url := range urls {
		id := u.ShortID(url)
		urlReturned, _ := store.GetURL(id) // MapStorage never returns non-nil error
		if urlReturned != url {
			t.Errorf("got [%s], want [%s]", urlReturned, url)
		}
	}
}

func getRandomFileName(t *testing.T) string {
	t.Helper()

	const length = 32
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
