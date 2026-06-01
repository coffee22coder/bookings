package config_test

import (
	"testing"

	"github.com/coffee22coder/bookings/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPlaceholder(t *testing.T) {
	require.True(t, true)
}

func TestLoad(t *testing.T) {
	locStr := "localhost"
	t.Setenv("DB_HOST", locStr)
	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, locStr, cfg.DBHost)
}

func TestLoad_MissingDBHost(t *testing.T) {
	t.Setenv("DB_HOST", "")
	cfg, err := config.Load()
	require.Nil(t, cfg)
	require.Error(t, err)
	require.ErrorContains(t, err, "DBHost")
}

func TestLoad_InvalidHTTPPort(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("HTTP_PORT", "0")
	_, err := config.Load()
	require.Error(t, err)
	require.ErrorContains(t, err, "HTTPPort")
}
