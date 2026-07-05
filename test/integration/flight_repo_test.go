//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/coffee22coder/bookings/internal/adapter/postgres"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/coffee22coder/bookings/internal/testutil"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestFlightRepo(t *testing.T) {
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

	repo := postgres.NewFlightRepo(pool)
	flights, err := repo.List(ctx, testutil.DepartureAirport, testutil.ArrivalAirport, testutil.Date, 2, 0)
	require.NoError(t, err)

	count, err := repo.CountSearch(ctx, testutil.DepartureAirport, testutil.ArrivalAirport, testutil.Date)
	require.NoError(t, err)

	require.GreaterOrEqual(t, count, len(flights))
}
