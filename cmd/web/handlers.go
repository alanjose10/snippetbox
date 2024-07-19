package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) homeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/pages/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) snippetViewGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprintf(w, "Displaying snippet with ID %d", id)
}

func (app *application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Disaply a form to create a new snippet")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "Test 1"
	content := "This is a test snippet - 1"
	expires := 1

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
