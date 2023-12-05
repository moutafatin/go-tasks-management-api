package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/api/v1/tasks", app.handleGetTasks)
	r.Get("/api/v1/tasks/{id}", app.handleGetTaskByID)

	r.Post("/api/v1/tasks", app.handleCreateTask)

	return r
}
