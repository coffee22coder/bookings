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
