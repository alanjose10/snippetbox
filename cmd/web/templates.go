package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/alanjose10/snippetbox/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func createTemplateCache() (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template)

	globMatches, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, file := range globMatches {
		name := filepath.Base(file)

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			file,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		// err := errors.New(fmt.Sprintf("the template %s does not exist", page))
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
