package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/shorty/internal/app/storage"
)

func GetHandler(store storage.Storage, serverName string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Software requirements specification:
		// ... Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
		// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
		// в HTTP-заголовке Location ...
		id := chi.URLParam(r, "id")
		if url, err := store.GetURL(id); err == nil {
			rw.Header().Set("Location", url)
			rw.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(rw, r)
	}
}

func PostHandler(store storage.Storage, serverName string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Software requirements specification:
		// ... Эндпоинт POST / принимает в теле запроса строку URL для сокращения
		// и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой
		// строки в теле ...
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Server failed to read the request's body", http.StatusInternalServerError)
			return
		}
		url := string(b)
		// TODO There should be some kind of URL validation
		if url == "" {
			http.Error(rw, "Bad request", http.StatusBadRequest)
			return
		}
		id := store.AddURL(url)
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, "http://%s/%s", serverName, id)
	}
}
