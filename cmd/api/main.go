package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error config", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	slog.Info("api starting", "http_port", cfg.HTTPPort)

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, *cfg)

	if err != nil {
		slog.Error("Error postgres", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	defer pool.Close()
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings.flights").Scan(&count)
	if err != nil {
		slog.Error("db query failed", "error", err)
		os.Exit(1)
	}
	slog.Info("db ok", "flights_count", count)

	os.Exit(0)
}
