package main

import (
	"errors"
	"net/http"
)

const defautNotFoundMessage = "the requested resource could not be found"

var ErrInvalidIdParam = errors.New("invalid id parameter")

func (app *application) logError(r *http.Request, err error) {
	var (
		uri    = r.URL.RequestURI()
		method = r.Method
	)
	app.logger.Error(err.Error(), "uri", uri, "method", method)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, status, envelope{"error": message}, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) fieldsErrorResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errs)
}
