package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alanjose10/snippetbox/internal/models"
	"github.com/alanjose10/snippetbox/internal/validators"
)

func (app *application) homeGet(w http.ResponseWriter, r *http.Request) {

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

	data := app.newTemplateData()

	data.Form = snippetCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.html", data)

}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validators.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validators.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validators.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validators.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validators.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
