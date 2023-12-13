package handlers

import (
	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/response"
)

type Config struct {
	Models data.Models
	Error  response.ErrorResponse
}

type Handlers struct {
	Tasks tasksHandler
}

func New(cfg Config) *Handlers {
	return &Handlers{
		Tasks: tasksHandler{
			models: cfg.Models,
			error:  cfg.Error,
		},
	}
}
