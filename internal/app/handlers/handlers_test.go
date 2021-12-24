package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"
)

func TestDefaultHandler_Post_Get_Valid(t *testing.T) {
	want_url := "http://localhost:8080/"
	tests := []struct {
		url       string
		id        string
		want_code int
	}{
		{
			url:       "http://example.com",
			id:        want_url + "1",
			want_code: 201,
		},
		{
			url:       "http://example.org",
			id:        want_url + "2",
			want_code: 201,
		},
		{
			url:       "http://local.test",
			id:        want_url + "3",
			want_code: 201,
		},
	}

	store := storage.NewMapStorage()

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, want_url, strings.NewReader(tt.url))
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()

			if resp.StatusCode != tt.want_code {
				t.Errorf("want status code [%d] but got [%d]", tt.want_code, resp.StatusCode)
			}

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("cannot read response body, should not see this normally")
			}

			if string(resBody) != tt.id {
				t.Errorf("want retirned id [%s] but got [%s]", tt.id, resBody)
			}
		})
	}

	for _, tt := range tests {
		t.Run("Get: "+tt.url, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.id, nil)
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()

			if resp.StatusCode != http.StatusTemporaryRedirect {
				t.Errorf("want status code [%d] but got [%d]",
					http.StatusTemporaryRedirect, resp.StatusCode)
			}

			location := resp.Header.Get("Location")

			if location != tt.url {
				t.Errorf("want retirned location header [%s] but got [%s]", location, tt.url)
			}
		})
	}
}

func TestDefaultHandler_Post_Not_Valid(t *testing.T) {
	want_url := "http://localhost:8080/"
	tests := []struct {
		url       string
		id        string
		want_code int
	}{
		{
			url:       "",
			id:        "Bad request\n",
			want_code: 400,
		},
	}
	store := storage.NewMapStorage()

	for _, tt := range tests {
		t.Run("Post: "+tt.id, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, want_url, strings.NewReader(tt.url))
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()

			if resp.StatusCode != tt.want_code {
				t.Errorf("want status code [%d] but got [%d]", tt.want_code, resp.StatusCode)
			}

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("cannot read response body, should not see this normally")
			}

			if string(resBody) != tt.id {
				t.Errorf("want retirned id [%s] but got [%s]", tt.id, resBody)
			}
		})
	}
}

func TestDefaultHandler_Get_Not_Valid(t *testing.T) {
	want_url := "http://localhost:8080/"
	tests := []struct {
		url       string
		id        string
		want_code int
	}{
		{
			url:       "",
			id:        want_url + "111111111",
			want_code: 404,
		},
		{
			url:       "",
			id:        want_url,
			want_code: 404,
		},
	}
	store := storage.NewMapStorage()
	for _, tt := range tests {
		t.Run("Get: "+tt.id, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.id, nil)
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()

			if resp.StatusCode != tt.want_code {
				t.Errorf("want status code [%d] but got [%d]",
					tt.want_code, resp.StatusCode)
			}

			location := resp.Header.Get("Location")

			if location != tt.url {
				t.Errorf("want retirned location header [%s] but got [%s]", location, tt.url)
			}
		})
	}
}
