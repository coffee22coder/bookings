package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/adapter/http/handler"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/coffee22coder/bookings/internal/port"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Services struct {
	Airports *service.AirportService
	Flights  *service.FlightService
}

type Server struct {
	cfg      *config.Config
	logger   *slog.Logger
	router   chi.Router
	db       port.DBPinger
	services *Services
}

func NewServer(cfg *config.Config, logger *slog.Logger, db port.DBPinger, services *Services) *Server {
	s := &Server{
		cfg:      cfg,
		logger:   logger,
		router:   chi.NewRouter(),
		db:       db,
		services: services,
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Recoverer)

	s.routes()

	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) routes() {
	airportHandler := handler.NewAirportHandler(s.services.Airports)
	flightHandler := handler.NewFlightHandler(s.services.Flights)

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s.logger.Info("health", "request_id", middleware.GetReqID(r.Context()))
		id := middleware.GetReqID(r.Context())
		w.Header().Set(middleware.RequestIDHeader, id)
		json.NewEncoder(w).Encode(dto.JSONResponse{Status: "ok"})
	})

	s.router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		if err := s.db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(dto.JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.JSONResponse{
			Status: "ok",
		})
	})

	s.router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	s.router.Get("/api/v1/airports", airportHandler.List)

	s.router.Get("/api/v1/airports/count", airportHandler.Count)

	s.router.Get("/api/v1/flights", flightHandler.List)

	s.router.Get("/api/v1/flights/count", flightHandler.CountSearch)
}
