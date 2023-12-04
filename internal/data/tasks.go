package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    string
	Status      string
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
