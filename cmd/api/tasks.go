package main

import (
	"net/http"
)

func (app *application) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.models.Tasks.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tasks": tasks}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
