package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello there"))
	})

	r.Get("/api/v1/tasks", handleGetTasks)

	return r
}
