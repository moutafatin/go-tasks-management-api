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
	r.Use(app.authenticate)

	r.Group(func(r chi.Router) {
		r.Use(app.requireActivatedUser)

		r.Get("/api/v1/tasks", app.handlers.Tasks.HandleGetTasks)
		r.Get("/api/v1/tasks/{id}", app.handlers.Tasks.HandleGetTaskByID)
		r.Delete("/api/v1/tasks/{id}", app.handlers.Tasks.HandleDeleteTask)
		r.Put("/api/v1/tasks/{id}", app.handlers.Tasks.HandleUpdateTask)
		r.Post("/api/v1/tasks", app.handlers.Tasks.HandleCreateTask)
	})

	r.Post("/api/v1/users", app.handleRegisterUser)
	r.Put("/api/v1/users/activated", app.handleActivateUser)

	r.Post("/api/v1/tokens/authentication", app.handleCreateAuthenticationToken)

	return r
}
