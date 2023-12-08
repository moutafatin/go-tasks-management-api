package main

import (
	"errors"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
)

func (app *application) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if data.ValidateUser(v, user); !v.Valid() {
		app.fieldsErrorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email already exists")
			app.fieldsErrorResponse(w, r, v.Errors)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
