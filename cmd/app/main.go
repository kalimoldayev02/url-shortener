package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/kalimoldayev02/shop/internal/repository/storage/postgres"
	"github.com/kalimoldayev02/shop/pkg/config"
	"github.com/kalimoldayev02/shop/pkg/lib/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shorter", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.GetStoragePath())
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	lastId, _ := storage.SaveUrl("long-url", "short-url")
	fmt.Println(lastId)
}

// настройки логирования
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
