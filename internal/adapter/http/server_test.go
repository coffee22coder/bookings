package http_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httpadapter "github.com/coffee22coder/bookings/internal/adapter/http"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/stretchr/testify/require"
)

func TestHealth_Returns200(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{HTTPPort: 8080}
	srv := httpadapter.NewServer(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"Status":"ok"}`, string(body))
}

func TestHealth_POST_Returns405(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{HTTPPort: 8080}
	srv := httpadapter.NewServer(cfg, logger)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHealth_ReturnsJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{HTTPPort: 8080}
	srv := httpadapter.NewServer(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	require.JSONEq(t, `{"Status":"ok"}`, rec.Body.String())
}
