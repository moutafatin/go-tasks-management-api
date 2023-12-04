package main

import (
	"errors"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
)

func (app *application) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.models.Tasks.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tasks": tasks}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIntParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.models.Tasks.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
