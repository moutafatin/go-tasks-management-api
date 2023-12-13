package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Tasks  tasksModel
	Users  usersModel
	Tokens tokensModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Tasks: tasksModel{
			DB: db,
		},
		Users: usersModel{
			DB: db,
		},
		Tokens: tokensModel{
			DB: db,
		},
	}
}
