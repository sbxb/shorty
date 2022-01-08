package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"

	"github.com/go-chi/chi/v5"
)

var cfg = config.DefaultConfig

func TestPostHandler_NotValidCases(t *testing.T) {
	wantCode := 400
	tests := []struct {
		url string
		id  string
	}{
		{url: ""},
	}

	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:      store,
		ServerName: cfg.FullServerName(),
	}
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.FullServerURL(), strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != wantCode {
				t.Errorf("want status code [%d],  got [%d]", wantCode, resp.StatusCode)
			}
		})
	}
}

func TestPostHandler_ValidCases(t *testing.T) {
	wantCode := 201
	tests := []struct {
		url string
		id  string
	}{
		{
			url: "http://example.com",
		},
		{
			url: "http://example.org",
		},
	}

	// Fill in test cases' ids
	for i, tt := range tests {
		tests[i].id = u.ShortID(tt.url)
	}

	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:      store,
		ServerName: cfg.FullServerName(),
	}
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.FullServerURL(), strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != wantCode {
				t.Errorf("want status code [%d],  got [%d]", wantCode, resp.StatusCode)
			}

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("cannot read response body, should not see this normally")
			}

			if string(resBody) != cfg.FullServerURL()+tt.id {
				t.Errorf("want returned id [%s],  got [%s]", cfg.FullServerURL()+tt.id, resBody)
			}
		})
	}
}

func TestGetHandler_NotValidCases(t *testing.T) {
	wantCode := 404
	tests := []struct {
		url string
		id  string
	}{
		{id: ""},
		{id: "NON_EXISTING_ID"},
	}
	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:      store,
		ServerName: cfg.FullServerName(),
	}
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.FullServerURL() + tt.id
		t.Run("Get: "+requestURL, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != wantCode {
				t.Errorf("want status code [%d],  got [%d]", wantCode, resp.StatusCode)
			}
		})
	}
}

func TestGetHandler_ValidCases(t *testing.T) {
	wantCode := 307
	tests := []struct {
		url string
		id  string
	}{
		{url: "http://example.com"},
		{url: "http://example.org"},
	}

	// Fill in test cases' ids
	for i, tt := range tests {
		tests[i].id = u.ShortID(tt.url)
	}

	// Prepare store
	store := storage.NewMapStorage()
	for _, tt := range tests {
		store.AddURL(tt.url, tt.id)
	}

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:      store,
		ServerName: cfg.FullServerName(),
	}
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.FullServerURL() + tt.id
		t.Run("Get: "+requestURL, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != wantCode {
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
