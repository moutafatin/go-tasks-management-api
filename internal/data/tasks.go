package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
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

func GetTaskStatus(status *string) TaskStatus {
	if status == nil {
		return taskStatusTodo
	}
	switch strings.ToLower(*status) {
	case "todo":
		return taskStatusTodo
	case "in_progress":
		return taskStatusInProgress
	case "done":
		return taskStatusDone
	case "":
		return ""
	default:
		return TaskStatus("invalid")
	}
}

func GetTaskPriority(priority *string) TaskPriority {
	if priority == nil {
		return TaskPriorityLow
	}
	switch strings.ToLower(*priority) {
	case "low":
		return TaskPriorityLow
	case "medium":
		return TaskPriorityMedium
	case "high":
		return TaskPriorityHigh
	case "":
		return ""
	default:
		return TaskPriority("invalid")
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

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(validator.NotEmpty(task.Title), "title", "title is required")
	v.Check(validator.NotEmpty(task.Description), "description", "description is required")
	v.Check(validator.NotEmpty(string(task.Priority)), "priority", "priority must not be empty if set")
	v.Check(validator.PremittedValues(task.Priority, []TaskPriority{TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh}), "priority", "Invalid priority value, must be one of `low`, `medium`, `high")
	v.Check(validator.NotEmpty(string(task.Status)), "status", "status must not be empty if set")
	v.Check(validator.PremittedValues(task.Status, []TaskStatus{taskStatusTodo, taskStatusInProgress, taskStatusDone}), "status", "Invalid status value, must be one of `todo`, `in_progress`, `done`")
}
