package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/prplx/go-sentry-tunnel/internal/config"
	"github.com/prplx/go-sentry-tunnel/internal/errors"
	"github.com/prplx/go-sentry-tunnel/internal/handlers"
	"github.com/prplx/go-sentry-tunnel/internal/lib/sl"
)

const (
	envDev  = "development"
	envProd = "production"
	envTest = "test"
)

func main() {
	config := config.MustLoad()
	log := setupLogger(config.Env)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tunnel", handlers.HandleTunnel)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		log.Error("Could not start server:", sl.Err(err))
		return
	case <-shutdown:
		log.Info("Starting shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Error("Could not gracefully shutdown the server:", sl.Err(err))
		}

		log.Info("Server gracefully stopped")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic(fmt.Errorf("%w: ENV", errors.ErrorEnvVariableRequired))
	}

	return log
}
