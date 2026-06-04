package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coffee22coder/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router chi.Router
}

func NewServer(cfg *config.Config, logger *slog.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		router: chi.NewRouter(),
	}

	s.routes()

	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.HTTPPort)
	s.logger.Info("listening", "addr", addr)
	return http.ListenAndServe(addr, s.router)
}

type Str struct {
	status string
}

func (s *Server) routes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct{ Status string }{Status: "ok"})

	})
}
