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
