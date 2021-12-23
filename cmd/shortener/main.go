package main

import (
	"net/http"

	"github.com/sbxb/shorty/internal/app/handlers"
)

const serverName = "localhost:8080"

var store = map[string]string{}

func main() {
	// http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "store: %v\n", store)
	// })

	http.Handle("/", handlers.DefaultHandler(store, serverName))

	http.ListenAndServe(serverName, nil)
}
