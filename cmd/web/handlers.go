package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alanjose10/snippetbox/internal/models"
)

func (app *application) homeGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.html", data)

}

func (app *application) snippetViewGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	s, err := app.snippetModel.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrSnippetNotFound) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Snippet = s

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Disaply a form to create a new snippet")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "Test 2"
	content := "This is a test snippet - 2"
	expires := 1

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
