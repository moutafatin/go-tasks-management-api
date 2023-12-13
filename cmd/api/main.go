package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/env"
	"github.com/moutafatin/go-tasks-management-api/internal/handlers"
	"github.com/moutafatin/go-tasks-management-api/internal/mailer"
	"github.com/moutafatin/go-tasks-management-api/internal/response"
	"github.com/subosito/gotenv"
)

type currentEnv string

func (e currentEnv) IsDevelopment() bool {
	return strings.Trim(string(e), " ") == "development"
}

type config struct {
	port        int
	environment string
	db          dbConfig
	limiter     struct {
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
	models   data.Models
	logger   *slog.Logger
	config   config
	mailer   mailer.Mailer
	error    response.ErrorResponse
	wg       sync.WaitGroup
	env      currentEnv
	handlers *handlers.Handlers
}

func main() {
	gotenv.Load()

	var cfg config

	flag.IntVar(&cfg.port, "port", env.GetInt("PORT", 4000), "TCP port to listen to")

	flag.StringVar(&cfg.environment, "environment", env.GetString("ENV", "development"), "environment development|staging|production")

	flag.StringVar(&cfg.db.dsn, "dsn", env.GetString("POSTGRES_URL", ""), "Postgres connection url")
	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-conns", env.GetInt("DB_MAX_OPEN_CONN", 25), "Pool maximum open connections")
	flag.StringVar(&cfg.db.maxIdleConnTime, "db-max-idle-time", env.GetString("DB_MAX_IDLE_TIME", "15m"), "Pool maximum idle connection time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", float64(env.GetInt("LIMITER_RPS", 2)), "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", env.GetInt("LIMITER_BURST", 4), "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", env.GetBool("LIMITER_ENABLED", false), "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", env.GetString("SMTP_HOST", ""), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", env.GetInt("SMTP_PORT", 0), "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", env.GetString("SMTP_USERNAME", ""), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", env.GetString("SMTP_PASSWORD", ""), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", env.GetString("SMTP_SENDER", ""), "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg.db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database")

	errorResponse := response.ErrorResponse{
		Logger: logger,
	}

	models := data.NewModels(db)
	app := &application{
		models: data.NewModels(db),
		logger: logger,
		config: cfg,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		error:  errorResponse, env: currentEnv(cfg.environment),
		handlers: handlers.New(handlers.Config{
			Error:  errorResponse,
			Models: models,
		}),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}
