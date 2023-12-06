package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.port),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutDownErrCh := make(chan error)

	go func() {
		quitCh := make(chan os.Signal, 1)

		signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

		s := <-quitCh
		app.logger.Info("Shutting down the server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		shutDownErrCh <- srv.Shutdown(ctx)
	}()

	app.logger.Info(fmt.Sprintf("server running on http://localhost%s", srv.Addr))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownErrCh
	if err != nil {
		return err
	}
	app.logger.Info("server stopped")
	return nil
}
