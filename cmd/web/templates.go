package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alanjose10/snippetbox/internal/models"
)

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 3:04 PM")
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

func (app *application) newTemplateData() (t templateData) {
	t = templateData{
		CurrentYear: time.Now().Year(),
	}
	return
}

func createTemplateCache() (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template)

	globMatches, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range globMatches {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(templateFunctions).ParseFiles("./ui/html/base.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]

	// initialize a new buffer
	buffer := new(bytes.Buffer)

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	err := ts.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	_, err = buffer.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
	}
}
