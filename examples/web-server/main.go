package main

import (
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	if err := StartServer(); err != nil {
		panic(err)
	}
}

// StartServer starts the webserver
func StartServer() error {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(MainView("Title", "Content")))
	return http.ListenAndServe(":8080", mux)
}
