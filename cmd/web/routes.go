package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a new server mux
	// and register the home function as handler for /
	mux := http.NewServeMux()

	fileServe := http.FileServer(http.Dir("./ui/static"))

	// File serve route
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServe))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is healthy"))
	})

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.homeGet))

	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetViewGet))

	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreateGet))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupGet))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginGet))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.

	return standard.Then(mux)
}
