package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
)

func (app *application) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Priority    *string `json:"priority"`
		Status      *string `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	log.Println("USER ID ", user.ID)
	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    data.GetTaskPriority(input.Priority),
		Status:      data.GetTaskStatus(input.Status),
		UserID:      user.ID,
	}
	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.faildErrorResponse(w, r, v.Errors)
		return
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
	user := app.contextGetUser(r)
	tasks, err := app.models.Tasks.GetAll(user.ID)
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
		app.badRequestResponse(w, r, ErrInvalidIdParam)
		return
	}

	user := app.contextGetUser(r)

	task, err := app.models.Tasks.GetByID(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r, "task not found")
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

func (app *application) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidIdParam)
		return
	}
	user := app.contextGetUser(r)
	err = app.models.Tasks.Delete(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r, "task not found")
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	// maybe return 201 no content, its depend
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "task deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidIdParam)
		return
	}

	user := app.contextGetUser(r)
	task, err := app.models.Tasks.GetByID(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r, "task not found")
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *string `json:"priority"`
		Status      *string `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Priority != nil {
		task.Priority = data.GetTaskPriority(input.Priority)
	}
	if input.Status != nil {
		task.Status = data.GetTaskStatus(input.Status)
	}

	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.faildErrorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tasks.Update(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "task updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
