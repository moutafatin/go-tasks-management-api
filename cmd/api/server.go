package main

import (
	"context"
	"crypto/tls"
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
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			// curves that have assembly implementation
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
		},
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
		err := srv.Shutdown(ctx)
		if err != nil {
			shutDownErrCh <- err
		}

		app.logger.Info("completing background tasks", "addr", srv.Addr)

		app.wg.Wait()
		shutDownErrCh <- nil
	}()

	app.logger.Info(app.getEnvBasedUrl())

	var err error

	if app.env.IsDevelopment() {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	}
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
