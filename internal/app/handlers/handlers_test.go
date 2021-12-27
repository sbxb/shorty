package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

var serverURL = "http://localhost:8080/"

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
	router.Post("/", handlers.PostHandler(store, "localhost:8080"))

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, serverURL, strings.NewReader(tt.url))
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
			id:  "1",
		},
		{
			url: "http://example.org",
			id:  "2",
		},
	}

	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	router.Post("/", handlers.PostHandler(store, "localhost:8080"))

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, serverURL, strings.NewReader(tt.url))
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

			if string(resBody) != serverURL+tt.id {
				t.Errorf("want returned id [%s],  got [%s]", serverURL+tt.id, resBody)
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
	router.Get("/{id}", handlers.GetHandler(store, "localhost:8080"))

	for _, tt := range tests {
		requestURL := serverURL + tt.id
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

	// Prepare store and fill in test cases' ids
	store := storage.NewMapStorage()
	for i, tt := range tests {
		tests[i].id = store.AddURL(tt.url)
	}

	router := chi.NewRouter()
	router.Get("/{id}", handlers.GetHandler(store, "localhost:8080"))

	for _, tt := range tests {
		requestURL := serverURL + tt.id
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
