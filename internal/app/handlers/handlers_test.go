package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/handlers"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"

	"github.com/go-chi/chi/v5"
)

// stackoverflow.com-recommended hack to parse testing flags before application ones
// prevents test failure with "flag provided but not defined: -test.testlogfile" error
var _ = func() bool {
	testing.Init()
	return true
}()

var cfg = config.New()

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

	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:  store,
		Config: cfg,
	}
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON", func(t *testing.T) {
			requestURL := cfg.BaseURL + "/api/shorten"
			requestBody, _ := json.Marshal(tt.requestObj)
			req := httptest.NewRequest(http.MethodPost, requestURL, bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != wantCode {
				t.Errorf("want status code [%d],  got [%d]", wantCode, resp.StatusCode)
			}

			// resBody, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	t.Fatalf("cannot read response body, should not see this normally")
			// }

			responseObj := u.URLResponse{}
			err := json.NewDecoder(resp.Body).Decode(&responseObj)
			if err != nil {
				t.Fatalf("cannot read response body, should not see this normally")
			}
			if !reflect.DeepEqual(responseObj, tt.wantResponseObj) {
				t.Errorf("want [%#v],  got [%#v]", tt.wantResponseObj, responseObj)
			}
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

	// Prepare empty store
	store := storage.NewMapStorage()

	router := chi.NewRouter()
	urlHandler := handlers.URLHandler{
		Store:  store,
		Config: cfg,
	}
	router.Post("/api/shorten", urlHandler.JSONPostHandler)

	for _, tt := range tests {
		t.Run("Post JSON", func(t *testing.T) {
			requestURL := cfg.BaseURL + "/api/shorten"
			req := httptest.NewRequest(http.MethodPost, requestURL, strings.NewReader(tt.body))
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
		Store:  store,
		Config: cfg,
	}
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(tt.url))
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
		Store:  store,
		Config: cfg,
	}
	router.Post("/", urlHandler.PostHandler)

	for _, tt := range tests {
		t.Run("Post: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, cfg.BaseURL+"/", strings.NewReader(tt.url))
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
			want := cfg.BaseURL + "/" + tt.id
			if string(resBody) != want {
				t.Errorf("want returned id [%s],  got [%s]", want, resBody)
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
		Store:  store,
		Config: cfg,
	}
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.BaseURL + "/" + tt.id
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
		Store:  store,
		Config: cfg,
	}
	router.Get("/{id}", urlHandler.GetHandler)

	for _, tt := range tests {
		requestURL := cfg.BaseURL + "/" + tt.id
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
