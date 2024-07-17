package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a new server mux
	// and register the home function as handler for /
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", homeGet)

	mux.HandleFunc("GET /snippet/view/{id}", snippetViewGet)

	mux.HandleFunc("GET /snippet/create", snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
