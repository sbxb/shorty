package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	store  storage.Storage
	config config.Config
}

func NewURLHandler(st storage.Storage, cfg config.Config) URLHandler {
	return URLHandler{
		store:  st,
		config: cfg,
	}
}

// GetHandler process GET /{id} request
// ... Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
// в HTTP-заголовке Location ...
func (uh URLHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := uh.store.GetURL(id)
	log.Printf("[%s] [%s]\n", id, url)
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

	id := u.ShortID(url)
	err = uh.store.AddURL(url, id)
	if err != nil {
		http.Error(w, "Server failed to store URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", uh.config.BaseURL, id)
}

// JSONPostHandler process POST /api/shorten request with JSON payload
// ... эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект
// {"url": "<some_url>"} и возвращающий в ответ объект {"result": "<shorten_url>"}
func (uh URLHandler) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	const ContentType = "application/json"
	var req u.URLRequest

	if r.Header.Get("Content-Type") != ContentType {
		http.Error(w, "Bad request: Content-Type should be "+ContentType, http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// is request an empty struct
	if req == (u.URLRequest{}) {
		http.Error(w, "Bad request: empty object received", http.StatusBadRequest)
		return
	}

	id := u.ShortID(req.URL)
	if err := uh.store.AddURL(req.URL, id); err != nil {
		http.Error(w, "Server failed to store URL", http.StatusInternalServerError)
		return
	}

	jr, err := json.Marshal(
		u.URLResponse{
			Result: fmt.Sprintf("%s/%s", uh.config.BaseURL, id),
		},
	)

	if err != nil {
		http.Error(w, "Server failed to process response result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(jr)
}
