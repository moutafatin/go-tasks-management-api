package response

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ErrorResponse struct {
	Logger *slog.Logger
}

func (e ErrorResponse) LogError(r *http.Request, err error) {
	var (
		uri    = r.URL.RequestURI()
		method = r.Method
		trace  = string(debug.Stack())
	)
	e.Logger.Error(err.Error(), "uri", uri, "method", method, "trace", trace)
}

func (e ErrorResponse) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := JSON(w, status, map[string]any{"error": message})
	if err != nil {
		e.LogError(r, err)
		w.WriteHeader(status)
	}
}

func (e ErrorResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.LogError(r, err)

	message := "the server encountered a problem and could not process your request"

	e.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (e ErrorResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e ErrorResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	e.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (e ErrorResponse) FaildErrorResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	e.ErrorResponse(w, r, http.StatusUnprocessableEntity, errs)
}

func (e ErrorResponse) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	e.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func (e ErrorResponse) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	e.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (e ErrorResponse) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	e.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (e ErrorResponse) UnAuthorizedResponse(w http.ResponseWriter, r *http.Request) {
	message := "you are not authorized to access this resource"
	e.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (e ErrorResponse) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	e.ErrorResponse(w, r, http.StatusForbidden, message)
}
