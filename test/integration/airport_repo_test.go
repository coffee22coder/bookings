//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestAirportRepo_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	require.NoError(t, err)

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, *cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	repo := postgres.New(pool)

	airports, err := repo.List(ctx, 5, 0)
	require.NoError(t, err)

	require.Len(t, airports, 5)
}

func TestAirportRepo_Count(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	require.NoError(t, err)

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, *cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	repo := postgres.New(pool)
	count, err := repo.Count(ctx)
	require.NoError(t, err)

	require.Greater(t, count, 0)
}
