package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"go-storage-api/internal/api"
	"go-storage-api/internal/config"
	"go-storage-api/internal/storage/local"
)

func main() {
	cfg := config.Load()

	level := parseLogLevel(cfg.LogLevel)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	store, err := local.New(cfg.Local.RootPath)
	if err != nil {
		log.Fatalf("create local storage backend: %v", err)
	}

	router := api.NewRouter(store, cfg.MaxUploadSize, logger)

	logger.Info("server started", "port", cfg.Port, "backend", cfg.StorageBackend)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
