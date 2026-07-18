package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/coffee22coder/bookings/internal/adapter/http"
	"github.com/coffee22coder/bookings/internal/adapter/idgen"
	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/coffee22coder/bookings/internal/service"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	p, err := postgres.NewPool(ctx, *cfg)

	if err != nil {
		logger.Error("Error datebase", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}
	defer p.Close()
	// var count int
	// if count, err = p.FligthsCount(ctx); err != nil {
	// 	logger.Error("Error datebase", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	// }
	// logger.Info("db ok", "flights_count", count)

	services := NewServices(p)

	server := httpadapter.NewServer(cfg, logger, p, services)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("HTTP server starting", "http_port", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("HTTP sever shutdown failed", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}

}

func NewServices(db *postgres.Pool) *httpadapter.Services {

	aeroportRepo := postgres.NewAeroportRepo(db)
	flightRepo := postgres.NewFlightRepo(db)

	genId := idgen.NewGenID(6)
	bookingRepo := postgres.NewBookingRepo(db, genId)

	return &httpadapter.Services{
		Airports: service.AirportServiceNew(aeroportRepo),
		Flights:  service.FlightServiceNew(flightRepo),
		Bookings: service.BookingServiceNew(bookingRepo),
	}
}
