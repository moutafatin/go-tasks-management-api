package main

import "net/http"

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
