package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const serverName = "localhost:8080"

var store = map[string]string{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
			// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
			// в HTTP-заголовке Location.
			id := strings.TrimLeft(r.URL.Path, "/")
			if url, ok := store[id]; ok {
				w.Header().Set("Location", url)
				w.WriteHeader(http.StatusTemporaryRedirect)
			} else {
				http.NotFound(w, r)
			}
		} else if r.Method == http.MethodPost {
			// Эндпоинт POST / принимает в теле запроса строку URL для сокращения
			// и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой
			// строки в теле.
			b, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Server failed to read the request's body", http.StatusInternalServerError)
				return
			}
			url := string(b)
			if url == "" {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			id := fmt.Sprintf("%d", len(store)+1)
			store[id] = url
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "http://%s/%s", serverName, id)
		}
	})

	http.ListenAndServe(serverName, nil)
}
