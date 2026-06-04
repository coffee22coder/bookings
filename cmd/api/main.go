package main

import (
	"context"
	"log/slog"
	"os"

	httpadapter "github.com/coffee22coder/bookings/internal/adapter/http"
	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Error config", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	logger.Info("api starting", "http_port", cfg.HTTPPort)

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, *cfg)

	if err != nil {
		logger.Error("Error postgres", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	defer pool.Close()
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings.flights").Scan(&count)
	if err != nil {
		logger.Error("db query failed", "error", err)
		os.Exit(1)
	}
	logger.Info("db ok", "flights_count", count)

	server := httpadapter.NewServer(cfg, logger)

	if err := server.Run(); err != nil {
		logger.Error("http server failed", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}

	os.Exit(0)
}
