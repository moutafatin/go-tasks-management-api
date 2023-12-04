package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/subosito/gotenv"
)

type dbConfig struct {
	dsn string
}

type config struct {
	port int
	db   dbConfig
}

type application struct {
	models data.Models
	logger *slog.Logger
}

func main() {
	gotenv.Load()

	var cfg config

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("POSTGRES_URL"), "Postgres connection url")
	flag.IntVar(&cfg.port, "port", getEnvInt("PORT"), "TCP port to listen to")

	logHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	db, err := openDb(cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database")

	app := &application{
		models: *data.NewModels(db),
		logger: logger,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.port),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logHandler, slog.LevelError),
	}

	logger.Info(fmt.Sprintf("server running on http://localhost%s", srv.Addr))
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDb(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
