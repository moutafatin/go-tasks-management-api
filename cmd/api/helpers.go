package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (app *application) readIntParam(r *http.Request, key string) (int, error) {
	stringId := chi.URLParam(r, key)

	id, err := strconv.Atoi(stringId)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}

func (app *application) getEnvBasedUrl() string {
	if app.env.IsDevelopment() {
		return fmt.Sprintf("http://localhost:%d", app.config.port)
	}

	return fmt.Sprintf("https://localhost:%d", app.config.port)
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("%v", err)
			}
		}()

		fn()
	}()
}
