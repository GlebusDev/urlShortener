package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/GlebusDev/urlShortener/internal/config"
	"github.com/GlebusDev/urlShortener/internal/httpServer/handlers/url/save"
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

	// router chi
	var router = chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(myMiddlewareLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, strg))
	router.Get("/{alias}", redirect.New(logger, strg))

	// server
	logger.Info("starting server", slog.String("address", config.Address))

	var server = &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.IdleTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server stopped")

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
