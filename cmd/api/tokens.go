package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/request"
	"github.com/moutafatin/go-tasks-management-api/internal/response"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
)

func (app *application) handleCreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := request.DecodeJSONStrict(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.faildErrorResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidCredentialsResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	matches, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, envelope{
		"authentication_token": token,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
