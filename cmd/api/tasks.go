package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
)

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []data.Task{
		{
			ID:          1,
			Title:       "Learn Go",
			Description: "Learn Go and build an API",
			Priority:    "High",
			Status:      "In Progress",
			DueDate:     time.Now().Add(time.Hour * 24 * 7),
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Learn React",
			Description: "Learn React and build a UI",
			Priority:    "High",
			Status:      "In Progress",
			DueDate:     time.Now().Add(time.Hour * 24 * 7),
			CreatedAt:   time.Now(),
		},
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
