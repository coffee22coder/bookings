package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/coffee22coder/bookings/internal/config"
	"github.com/coffee22coder/bookings/internal/port"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Services struct {
	Airports *service.AirportService
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

type JSONResponse struct {
	Status     string `json:"status"`
	ErrMessage string `json:"error,omitempty"`
}

func (s *Server) routes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s.logger.Info("health", "request_id", middleware.GetReqID(r.Context()))
		id := middleware.GetReqID(r.Context())
		w.Header().Set(middleware.RequestIDHeader, id)
		json.NewEncoder(w).Encode(JSONResponse{Status: "ok"})
	})

	s.router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		if err := s.db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSONResponse{
			Status: "ok",
		})
	})

	s.router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	s.router.Get("/api/v1/airports", func(w http.ResponseWriter, r *http.Request) {
		limit := 5
		offset := 0
		var err error
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(JSONResponse{
					Status:     "down",
					ErrMessage: err.Error(),
				})
				return
			}
		}
		offsetStr := r.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(JSONResponse{
					Status:     "down",
					ErrMessage: err.Error(),
				})
				return
			}
		}

		airports, err := s.services.Airports.List(r.Context(), limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(airports)
	})

	s.router.Get("/api/v1/airports/count", func(w http.ResponseWriter, r *http.Request) {
		count, err := s.services.Airports.Count(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(count)
	})
}
