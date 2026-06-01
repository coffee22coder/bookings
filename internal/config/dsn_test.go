package config_test

import (
	"testing"

	"github.com/coffee22coder/bookings/internal/config"
	"github.com/stretchr/testify/require"
)

func TestDSN(t *testing.T) {
	dsnMock := "postgres://test:test@localhost:1234/db?sslmode=disable"
	t.Setenv("DB_USER", "test")
	t.Setenv("DB_PASSWORD", "test")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "1234")
	t.Setenv("DB_NAME", "db")
	t.Setenv("DB_SSLMODE", "disable")
	cfg, err := config.Load()
	require.NoError(t, err)
	dsn := config.DSN(*cfg)
	require.NotEmpty(t, dsn)
	require.Equal(t, dsnMock, dsn)
}
