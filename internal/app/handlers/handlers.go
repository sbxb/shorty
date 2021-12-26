package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sbxb/shorty/internal/app/storage"
)

func DefaultHandler(store storage.Storage, serverName string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var id string
		if r.Method == http.MethodGet {
			// Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
			// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
			// в HTTP-заголовке Location.
			id = strings.TrimLeft(r.URL.Path, "/")
			if url, err := store.GetURL(id); err == nil {
				rw.Header().Set("Location", url)
				rw.WriteHeader(http.StatusTemporaryRedirect)
			} else {
				http.NotFound(rw, r)
			}
		} else if r.Method == http.MethodPost {
			// Эндпоинт POST / принимает в теле запроса строку URL для сокращения
			// и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой
			// строки в теле.
			b, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(rw, "Server failed to read the request's body", http.StatusInternalServerError)
				return
			}
			url := string(b)
			if url == "" {
				http.Error(rw, "Bad request", http.StatusBadRequest)
				return
			}
			id = store.AddURL(url)
			rw.WriteHeader(http.StatusCreated)
			fmt.Fprintf(rw, "http://%s/%s", serverName, id)
		} else {
			http.Error(rw, "Bad request", http.StatusBadRequest)
		}
	}
}
