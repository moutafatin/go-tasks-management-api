package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
	"github.com/moutafatin/go-tasks-management-api/internal/mailer"
	"github.com/subosito/gotenv"
)

type env struct {
	environment string
}

func (e env) IsDevelopment() bool {
	return strings.Trim(e.environment, " ") == "development"
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
	models data.Models
	logger *slog.Logger
	config config
	mailer mailer.Mailer
	wg     sync.WaitGroup
	env    env
}

func main() {
	gotenv.Load()

	var cfg config

	flag.IntVar(&cfg.port, "port", getEnvInt("PORT"), "TCP port to listen to")

	flag.StringVar(&cfg.environment, "environment", os.Getenv("ENV"), "environment development|staging|production")

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("POSTGRES_URL"), "Postgres connection url")
	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-conns", getEnvInt("DB_MAX_OPEN_CONN"), "Pool maximum open connections")
	flag.StringVar(&cfg.db.maxIdleConnTime, "db-max-idle-time", os.Getenv("DB_MAX_IDLE_TIME"), "Pool maximum idle connection time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", float64(getEnvInt("LIMITER_RPS")), "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", getEnvInt("LIMITER_BURST"), "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", getEnvBool("LIMITER_ENABLED"), "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", getEnvInt("SMTP_PORT"), "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg.db)
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
		env: env{
			environment: cfg.environment,
		},
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}
