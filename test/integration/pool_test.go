//go:build integration

package integration_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/coffee22coder/bookings/internal/adapter/http"
	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestPool_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	require.NoError(t, err)

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, *cfg)
	require.NoError(t, err)

	require.NotNil(t, pool)
	pool.Close()
}

func TestReady_WithRealPool_Returns200(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_ = godotenv.Load("../../.env")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg, err := config.Load()
	require.NoError(t, err)
	ctx := context.Background()
	p, err := postgres.NewPool(ctx, *cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		p.Close()
	})
	server := httpadapter.NewServer(cfg, logger, p)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}
