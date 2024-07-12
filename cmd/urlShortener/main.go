package main

import (
	"compress/bzip2"
	"fmt"
	"log/slog"
	"os"

	"github.com/GlebusDev/urlShortener/internal/config"
	myMiddlewareLogger "github.com/GlebusDev/urlShortener/internal/httpServer/middleware/logger"
	"github.com/GlebusDev/urlShortener/internal/lib/logger/sl"
	"github.com/GlebusDev/urlShortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	_ = strg
		// router chi
	var router = chi.NewRouter()


	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(myMiddlewareLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
		
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
