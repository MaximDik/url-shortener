package main

import (
	"log/slog"
	"os"

	"github.com/MaximDik/url-shortener/internal/config"
	"github.com/MaximDik/url-shortener/internal/lib/logger/sl"
	"github.com/MaximDik/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://yandex.com", "yandex")
	if err != nil {
		log.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("saved url", slog.Int64("id", id))

	id, err = storage.SaveURL("https://yandex.com", "yandex")
	if err != nil {
		log.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init logger:slog
	// TODO: init storage: sqlite
	// TODO: init router: chi, "chirender"
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
