package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type TaskStatus string

const (
	taskStatusTodo       TaskStatus = "TODO"
	taskStatusInProgress TaskStatus = "IN_PROGRESS"
	taskStatusDone       TaskStatus = "DONE"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
)

func GetTaskStatus(status string) TaskStatus {
	switch strings.ToLower(status) {
	case "todo":
		return taskStatusTodo
	case "in_progress":
		return taskStatusInProgress
	case "done":
		return taskStatusDone
	default:
		return ""
	}
}

func GetTaskPriority(priority string) TaskPriority {
	switch strings.ToLower(priority) {
	case "low":
		return TaskPriorityLow
	case "medium":
		return TaskPriorityMedium
	case "high":
		return TaskPriorityHigh
	default:
		return ""
	}
}

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    TaskPriority
	Status      TaskStatus
	CreatedAt   time.Time `db:"created_at"`
}

type tasksModel struct {
	DB *pgxpool.Pool
}

func (t *tasksModel) GetAll() ([]*Task, error) {
	stmt := `SELECT id, title, description, priority, status, created_at FROM tasks`

	rows, err := t.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Task])
}

func (t *tasksModel) GetByID(id int) (*Task, error) {
	stmt := `SELECT id, title, description, priority, status, created_at FROM tasks WHERE id = $1`

	row := t.DB.QueryRow(context.Background(), stmt, id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (t *tasksModel) Insert(task *Task) error {
	stmt := `INSERT INTO tasks (title, description, priority, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	args := []any{task.Title, task.Description, task.Priority, task.Status}

	return t.DB.QueryRow(context.Background(), stmt, args...).Scan(&task.ID, &task.CreatedAt)
}
