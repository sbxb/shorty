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
	wantURL := "http://localhost:8080/"
	tests := []struct {
		url      string
		id       string
		wantCode int
	}{
		{
			url:      "http://example.com",
			id:       wantURL + "1",
			wantCode: 201,
		},
		{
			url:      "http://example.org",
			id:       wantURL + "2",
			wantCode: 201,
		},
		{
			url:      "http://local.test",
			id:       wantURL + "3",
			wantCode: 201,
		},
	}

	store := storage.NewMapStorage()

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, wantURL, strings.NewReader(tt.url))
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("want status code [%d] but got [%d]", tt.wantCode, resp.StatusCode)
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
			defer resp.Body.Close()

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
	wantURL := "http://localhost:8080/"
	tests := []struct {
		url      string
		id       string
		wantCode int
	}{
		{
			url:      "",
			id:       "Bad request\n",
			wantCode: 400,
		},
	}
	store := storage.NewMapStorage()

	for _, tt := range tests {
		t.Run("Post: "+tt.id, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, wantURL, strings.NewReader(tt.url))
			w := httptest.NewRecorder()

			h := handlers.DefaultHandler(store, "localhost:8080")
			h.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("want status code [%d] but got [%d]", tt.wantCode, resp.StatusCode)
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
	wantURL := "http://localhost:8080/"
	tests := []struct {
		url      string
		id       string
		wantCode int
	}{
		{
			url:      "",
			id:       wantURL + "111111111",
			wantCode: 404,
		},
		{
			url:      "",
			id:       wantURL,
			wantCode: 404,
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
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("want status code [%d] but got [%d]",
					tt.wantCode, resp.StatusCode)
			}

			location := resp.Header.Get("Location")

			if location != tt.url {
				t.Errorf("want retirned location header [%s] but got [%s]", location, tt.url)
			}
		})
	}
}
