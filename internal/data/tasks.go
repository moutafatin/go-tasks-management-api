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
	UserID      int       `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type tasksModel struct {
	DB *pgxpool.Pool
}

func (t *tasksModel) GetAll(userID int) ([]*Task, error) {
	stmt := `SELECT id, title, description, priority, status,user_id, created_at FROM tasks WHERE user_id = $1`

	rows, err := t.DB.Query(context.Background(), stmt, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Task])
}

func (t *tasksModel) GetByID(id, userID int) (*Task, error) {
	stmt := `SELECT id, title, description, priority, status, user_id, created_at FROM tasks WHERE id = $1 AND user_id = $2`

	row := t.DB.QueryRow(context.Background(), stmt, id, userID)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.UserID, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (t *tasksModel) Insert(task *Task) error {
	stmt := `INSERT INTO tasks (title, description, priority, status,user_id) VALUES ($1, $2, $3, $4,$5) RETURNING id, created_at`

	args := []any{task.Title, task.Description, task.Priority, task.Status, task.UserID}

	return t.DB.QueryRow(context.Background(), stmt, args...).Scan(&task.ID, &task.CreatedAt)
}

func (t *tasksModel) Delete(id, userID int) error {
	stmt := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	res, err := t.DB.Exec(context.Background(), stmt, id, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return ErrRecordNotFound
	}
	return nil
}

func (t *tasksModel) Update(task *Task) error {
	stmt := `UPDATE tasks SET title = $1, description = $2, priority = $3, status = $4 WHERE id = $5 and user_id = $6`
	args := []any{task.Title, task.Description, task.Priority, task.Status, task.ID, task.UserID}

	_, err := t.DB.Exec(context.Background(), stmt, args...)

	return err
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(validator.NotEmpty(task.Title), "title", "title is required")
	v.Check(validator.NotEmpty(task.Description), "description", "description is required")
	v.Check(validator.NotEmpty(string(task.Priority)), "priority", "priority must not be empty if set")
	v.Check(validator.PremittedValues(task.Priority, []TaskPriority{TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh}), "priority", "Invalid priority value, must be one of `low`, `medium`, `high")
	v.Check(validator.NotEmpty(string(task.Status)), "status", "status must not be empty if set")
	v.Check(validator.PremittedValues(task.Status, []TaskStatus{taskStatusTodo, taskStatusInProgress, taskStatusDone}), "status", "Invalid status value, must be one of `todo`, `in_progress`, `done`")
	if task.UserID < 1 {
		panic("invalid operation,task can't exist without a user")
	}
}
