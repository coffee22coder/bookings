package postgres_test

import (
	"context"
	"testing"

	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewPool_InvalidHost(t *testing.T) {
	locStr := "local"
	t.Setenv("DB_HOST", locStr)
	cfg, err := config.Load()
	require.NoError(t, err)
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, *cfg)
	require.Error(t, err)
	require.ErrorContains(t, err, "ping db")
	require.Nil(t, pool)
}

func TestNewPool_NoPanicToClose(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "avia")
	t.Setenv("DB_PASSWORD", "avia")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "demo")
	cfg, err := config.Load()
	require.NoError(t, err)
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, *cfg)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	pool.Close()
	require.NotPanics(t, func() { pool.Close() })
}
