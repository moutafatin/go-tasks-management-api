package data

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    string
	Status      string
	DueDate     time.Time
	CreatedAt   time.Time
}
