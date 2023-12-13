package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbConfig struct {
	dsn             string
	maxOpenConn     int
	maxIdleConnTime string
}

func openDb(cfg dbConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.dsn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(cfg.maxOpenConn)

	duration, err := time.ParseDuration(cfg.maxIdleConnTime)
	if err != nil {
		return nil, err
	}

	config.MaxConnIdleTime = duration

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
