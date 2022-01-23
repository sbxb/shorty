package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg config.Config

var _ = func() bool {
	// stackoverflow.com-recommended hack to parse testing flags before
	// application ones prevents test failure with an error:
	// "flag provided but not defined: -test.testlogfile"
	testing.Init()

	var err error
	if cfg, err = config.New(); err != nil {
		log.Fatal(err)
	}
	return true
}()

func TestJSONPostHandler_ValidCases(t *testing.T) {
	wantCode := 201
	tests := []struct {
		url             string
		requestObj      u.URLRequest
		wantResponseObj u.URLResponse
	}{
		{
			url:             "http://example.com",
			requestObj:      u.URLRequest{},
			wantResponseObj: u.URLResponse{},
		},
		{
			url:             "http://example.org",
			requestObj:      u.URLRequest{},
			wantResponseObj: u.URLResponse{},
		},
	}

	// Fill in test cases' request and response objects
	for i, tt := range tests {
		tests[i].requestObj.URL = tt.url
		tests[i].wantResponseObj.Result = fmt.Sprintf("%s/%s", cfg.BaseURL, u.ShortID(tt.url))
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON", func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.requestObj)
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)

			responseObj := u.URLResponse{}

			err := json.NewDecoder(resp.Body).Decode(&responseObj)
			require.NoError(t, err, "cannot read response body, should not see this normally")

			assert.Equal(t, responseObj, tt.wantResponseObj)
		})
	}
}

func TestJSONPostHandler_NotValidCases(t *testing.T) {
	wantCode := 400
	tests := []struct {
		body string
	}{
		{""},
		{" "},
		{"abc"},
		{"{"},
		{"{}"},
		{`{"key": "value"}`},
		{`{"url": "http://example.com", "key": "value"}`}, // extra field
		{`{"url": ""}`}, // empty url
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)
		})
	}
}

func TestPostHandler_NotValidCases(t *testing.T) {
	wantCode := 400
	tests := []struct {
		url string
		id  string
	}{
		{url: ""},
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)
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

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)

			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "cannot read response body, should not see this normally")

			want := cfg.BaseURL + "/" + tt.id
			assert.Equal(t, string(resBody), want)
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

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.BaseURL + "/" + tt.id
		t.Run("Get: "+requestURL, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)
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

	store, _ := storage.NewMapStorage()
	for _, tt := range tests {
		store.AddURL(tt.url, tt.id)
	}

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.BaseURL + "/" + tt.id
		t.Run("Get: "+requestURL, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)

			location := resp.Header.Get("Location")

			assert.Equal(t, location, tt.url, "location header")
		})
	}
}
