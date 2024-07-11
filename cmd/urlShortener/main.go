package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/GlebusDev/urlShortener/internal/config"
	"github.com/GlebusDev/urlShortener/internal/lib/logger/sl"
	"github.com/GlebusDev/urlShortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// config cleanenv
	var config = config.MustLoad()

	fmt.Print(config)

	// logger slog
	var logger = setupLogger(config.Env)

	logger.Info("starting app ", slog.String("env", config.Env))
	logger.Debug("Debug messages are enabled")

	// storage sqlite
	strg, err := sqlite.New(config.StoragePath)

	if err != nil {
		logger.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	outUrl, err := strg.GetURL("shot")

	if err != nil {
		logger.Error("Faild to get url", sl.Err(err))
	}

	logger.Info("URL: ", outUrl)
	// router chi

	// server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
