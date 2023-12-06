package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(app.rateLimit)

	r.Get("/api/v1/tasks", app.handleGetTasks)
	r.Get("/api/v1/tasks/{id}", app.handleGetTaskByID)
	r.Delete("/api/v1/tasks/{id}", app.handleDeleteTask)
	r.Put("/api/v1/tasks/{id}", app.handleUpdateTask)

	r.Post("/api/v1/tasks", app.handleCreateTask)

	return r
}
