package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a new server mux
	// and register the home function as handler for /
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("Starting server on :3000")

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)

}

func home(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello from Snippetbox")
}
