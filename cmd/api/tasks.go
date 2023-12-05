package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
)

// TODO: finish this handler
func (app *application) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string
		Description string
		Priority    string
		Status      string
	}

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: validate input

	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    data.GetTaskPriority(input.Priority),
		Status:      data.GetTaskStatus(input.Status),
	}
	err = app.models.Tasks.Insert(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.Header().Add("location", fmt.Sprint("api/v1/tasks/", task.ID))
	app.writeJSON(w, http.StatusCreated, envelope{"task": task}, w.Header())
}

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
	id, err := app.readIntParam(r, "id")
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
