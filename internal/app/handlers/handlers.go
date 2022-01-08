package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	Store      storage.Storage
	ServerName string
}

// GetHandler process GET /{id} request
// ... Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
// в HTTP-заголовке Location ...
func (uh URLHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, err := uh.Store.GetURL(id)
	if err != nil {
		http.Error(w, "Server failed to process URL", http.StatusInternalServerError)
		return
	}
	if url == "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PostHandler process POST / request
// ... Эндпоинт POST / принимает в теле запроса строку URL для сокращения
// и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой
// строки в теле ...
func (uh URLHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server failed to read the request's body", http.StatusInternalServerError)
		return
	}

	url := string(b)
	// TODO There should be some kind of URL validation
	if url == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := u.ShortId(url)
	err = uh.Store.AddURL(url, id)
	if err != nil {
		http.Error(w, "Server failed to store URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%s", uh.ServerName, id)
}
