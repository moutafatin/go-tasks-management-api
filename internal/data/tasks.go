package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
