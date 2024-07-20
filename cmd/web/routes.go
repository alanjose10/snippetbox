package main

import "net/http"

func (app *application) routes() http.Handler {
	// Create a new server mux
	// and register the home function as handler for /
	mux := http.NewServeMux()

	fileServe := http.FileServer(http.Dir("./ui/static"))

	// File serve route
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServe))

	mux.HandleFunc("GET /{$}", app.homeGet)

	mux.HandleFunc("GET /snippet/view/{id}", app.snippetViewGet)

	mux.HandleFunc("GET /snippet/create", app.snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return app.logRequest(commonHeaders(mux))
}
