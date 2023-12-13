package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/ctx"
	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/request"
	"github.com/moutafatin/go-tasks-management-api/internal/response"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
)

type tasksHandler struct {
	models data.Models
	error  response.ErrorResponse
}

func (t tasksHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Priority    *string `json:"priority"`
		Status      *string `json:"status"`
	}

	err := request.DecodeJSONStrict(w, r, &input)
	if err != nil {
		t.error.BadRequestResponse(w, r, err)
		return
	}

	user := ctx.ContextGetUser(r)

	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    data.GetTaskPriority(input.Priority),
		Status:      data.GetTaskStatus(input.Status),
		UserID:      user.ID,
	}
	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		t.error.FaildErrorResponse(w, r, v.Errors)
		return
	}
	err = t.models.Tasks.Insert(task)
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
		return
	}
	w.Header().Add("location", fmt.Sprint("api/v1/tasks/", task.ID))
	err = response.JSONWithHeaders(w, http.StatusCreated, response.Envelope{"task": task}, w.Header())
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
	}
}

func (t tasksHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	user := ctx.ContextGetUser(r)
	tasks, err := t.models.Tasks.GetAll(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = response.JSON(w, http.StatusOK, response.Envelope{"tasks": tasks})
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
	}
}

func (t tasksHandler) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "id")
	if err != nil {
		t.error.BadRequestResponse(w, r, ErrInvalidIdParam)
		return
	}

	user := ctx.ContextGetUser(r)

	task, err := t.models.Tasks.GetByID(id, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			t.error.NotFoundResponse(w, r, "task not found")
			return
		}

		t.error.ServerErrorResponse(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, response.Envelope{"task": task})
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
	}
}

func (t tasksHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "id")
	if err != nil {
		t.error.BadRequestResponse(w, r, ErrInvalidIdParam)
		return
	}
	user := ctx.ContextGetUser(r)
	err = t.models.Tasks.Delete(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			t.error.NotFoundResponse(w, r, "task not found")
		default:
			t.error.ServerErrorResponse(w, r, err)
		}
		return
	}

	// maybe return 201 no content, its depend
	err = response.JSON(w, http.StatusOK, response.Envelope{"message": "task deleted successfully"})
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
	}
}

func (t tasksHandler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "id")
	if err != nil {
		t.error.BadRequestResponse(w, r, ErrInvalidIdParam)
		return
	}

	user := ctx.ContextGetUser(r)

	task, err := t.models.Tasks.GetByID(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			t.error.NotFoundResponse(w, r, "task not found")
		default:
			t.error.ServerErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *string `json:"priority"`
		Status      *string `json:"status"`
	}

	err = request.DecodeJSONStrict(w, r, &input)
	if err != nil {
		t.error.BadRequestResponse(w, r, err)
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
		t.error.FaildErrorResponse(w, r, v.Errors)
		return
	}

	err = t.models.Tasks.Update(task)
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, response.Envelope{"message": "task updated successfully"})
	if err != nil {
		t.error.ServerErrorResponse(w, r, err)
		return
	}
}
