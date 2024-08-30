package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-sentry-tunnel/internal/config"
	"go-sentry-tunnel/internal/handlers"
	"go-sentry-tunnel/internal/lib/sl"
)

const (
	envProd = "production"
)

func main() {
	config := config.MustLoad()
	log := setupLogger(config.Env)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tunnel", handlers.HandleTunnel(log, config))

	log = log.With("op", "cmd/api/main")

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}
	serverErrors := make(chan error, 1)

	go func() {
		log.Info("Starting server on port: " + config.Port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		log.Error("Could not start server:", sl.Err(err))
		return
	case <-shutdown:
		log.Info("Starting shutdown")

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

	if env == envProd {
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	} else {
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
