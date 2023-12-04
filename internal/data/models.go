package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Tasks tasksModel
}

func NewModels(db *pgxpool.Pool) *Models {
	return &Models{
		Tasks: tasksModel{
			DB: db,
		},
	}
}
