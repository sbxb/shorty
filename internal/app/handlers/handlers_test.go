package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestJSONBatchPostHandler_NotValidCases(t *testing.T) {
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
		{`{"correlation_id": "123", "original_url": "http://example.com"}`}, // not within an array
		{"["},
		{`[{"key": "value"}]`},
		{`[{"correlation_id": "123", "original_url": "http://example.com", "extra": "field"}]`}, // extra field
		{`[]`},                          // empty array
		{`[{"correlation_id": "123"}]`}, // no field
		{`[{"correlation_id": "", "original_url": "http://example.com"}]`}, // empty field
		//{`[{"correlation_id": "123", "original_url": "http://example.com"}]`}, // valid, commented only
	}
	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten/batch", urlHandler.JSONBatchPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten/batch", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			//resBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, resp.StatusCode, wantCode)
			//t.Log(string(resBody))
		})
	}
}

func TestJSONBatchPostHandler_ValidCases(t *testing.T) {
	wantCode := 201
	requestObj := []u.BatchURLRequestEntry{
		{
			CorrelationID: "123",
			OriginalURL:   "http://example.com",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "http://example.org",
		},
	}

	wantResponseObj := []u.BatchURLEntry{
		{
			CorrelationID: "123",
			OriginalURL:   "",
			ShortURL:      cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3",
		},
		{
			CorrelationID: "456",
			OriginalURL:   "",
			ShortURL:      cfg.BaseURL + "/6EH6vwAy9dOyyNbopTS6M4",
		},
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten/batch", urlHandler.JSONBatchPostHandler)

	requestBody, _ := json.Marshal(requestObj)
	req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten/batch", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, wantCode)

	responseObj := []u.BatchURLEntry{}

	err := json.NewDecoder(resp.Body).Decode(&responseObj)
	require.NoError(t, err, "cannot read response body, should not see this normally")

	assert.Equal(t, responseObj, wantResponseObj)
}

func TestJSONPostHandler_ValidCases(t *testing.T) {
	wantCode := 201
	tests := []struct {
		url              string
		wantResult       string
		buildInputOutput func(string, string) (u.URLRequest, u.URLResponse)
	}{
		{
			url:              "http://example.com",
			wantResult:       cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3",
			buildInputOutput: getRequestResponse,
		},
		{
			url:              "http://example.org",
			wantResult:       cfg.BaseURL + "/6EH6vwAy9dOyyNbopTS6M4",
			buildInputOutput: getRequestResponse,
		},
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON "+tt.url, func(t *testing.T) {
			requestObj, wantResponseObj := tt.buildInputOutput(tt.url, tt.wantResult)
			requestBody, _ := json.Marshal(requestObj)
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

			assert.Equal(t, responseObj, wantResponseObj)
		})
	}
}

func TestJSONPostHandler_InputRepeated(t *testing.T) {
	wantCode := 409
	requestObj := u.URLRequest{URL: "http://example.com"}
	wantResponseObj := u.URLResponse{Result: cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3"}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	requestBody, _ := json.Marshal(requestObj)

	// send the request
	req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ... and the same request once more
	req = httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, wantCode)

	responseObj := u.URLResponse{}

	err := json.NewDecoder(resp.Body).Decode(&responseObj)
	require.NoError(t, err, "cannot read response body, should not see this normally")
	assert.Equal(t, responseObj, wantResponseObj)
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
		url  string
		want string
	}{
		{
			url:  "http://example.com",
			want: cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3",
		},
		{
			url:  "http://example.org",
			want: cfg.BaseURL + "/6EH6vwAy9dOyyNbopTS6M4",
		},
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

			assert.Equal(t, string(resBody), tt.want)
		})
	}
}

func TestPostHandler_InputRepeated(t *testing.T) {
	wantCode := 409
	url := "http://example.com"
	want := cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3"

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Post("/", urlHandler.PostHandler)

	// send the request
	req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(url))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ... and the same request once more
	req = httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(url))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, wantCode)

	resBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "cannot read response body, should not see this normally")

	assert.Equal(t, string(resBody), want)
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
		reqURL  string
		wantURL string
	}{
		{
			reqURL:  cfg.BaseURL + "/5agFZWrIb6Ej21QvYUNBL3",
			wantURL: "http://example.com",
		},
		{
			reqURL:  cfg.BaseURL + "/6EH6vwAy9dOyyNbopTS6M4",
			wantURL: "http://example.org",
		},
	}

	store, _ := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.NewURLHandler(store, cfg)
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		store.AddURL(context.Background(), u.URLEntry{
			ShortURL:    tt.reqURL[strings.LastIndex(tt.reqURL, "/")+1:],
			OriginalURL: tt.wantURL,
		}, "")
		t.Run("Get: "+tt.reqURL, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.reqURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, wantCode)

			location := resp.Header.Get("Location")

			assert.Equal(t, location, tt.wantURL, "location header")
		})
	}
}

func getRequestResponse(url string, result string) (u.URLRequest, u.URLResponse) {
	return u.URLRequest{
			URL: url,
		},
		u.URLResponse{
			Result: result,
		}
}
