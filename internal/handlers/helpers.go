package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var ErrInvalidIdParam = errors.New("invalid id parameter")

func readIntParam(r *http.Request, key string) (int, error) {
	stringId := chi.URLParam(r, key)

	id, err := strconv.Atoi(stringId)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}
