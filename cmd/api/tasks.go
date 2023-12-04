package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.models.Tasks.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(js)
}
