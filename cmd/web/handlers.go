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

	data := app.newTemplateData(r)
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

	data := app.newTemplateData(r)
	data.Snippet = s

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.html", data)

}

type snippetCreateForm struct {
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validators.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validators.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validators.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validators.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Auth handlers

type userSignupForm struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

func (app *application) userSignupGet(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)

}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validators.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	id, err := app.userModel.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.logger.Debug(fmt.Sprintf("User created. Id: %d ", id))
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

func (app *application) userLoginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	form.CheckField(validators.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validators.Matches(form.Email, validators.EmailRX), "email", "Invalid email address")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	userId, err := app.userModel.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.logger.Debug(err.Error())
			form.AddNonFieldError("Invalid email/password. Please try again.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return

		} else if errors.Is(err, models.ErrUserDoesNotExist) {
			app.logger.Debug(err.Error())
			form.AddNonFieldError("User does not exist. Please signup.")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else {
			app.serverError(w, r, err)
		}
	}

	app.logger.Debug(fmt.Sprintf("Successfully authenticated user %d\n", userId))

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", userId)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
