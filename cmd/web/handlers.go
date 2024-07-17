package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) homeGet(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/pages/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) snippetViewGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Displaying snippet with ID %d", id)
}

func (app *application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Disaply a form to create a new snippet")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "New post created")

}
