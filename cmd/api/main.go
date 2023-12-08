package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/mailer"
	"github.com/subosito/gotenv"
)

type dbConfig struct {
	dsn string
}

type config struct {
	port    int
	db      dbConfig
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	models data.Models
	logger *slog.Logger
	config config
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	gotenv.Load()

	var cfg config

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("POSTGRES_URL"), "Postgres connection url")
	flag.IntVar(&cfg.port, "port", getEnvInt("PORT"), "TCP port to listen to")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "573515e3e82f45", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "ae4eb47eb4801d", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Taskio <no-reply@moutafatin.dev>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
		config: cfg,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
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
