package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestErrNotFound(t *testing.T) {
	err := fmt.Errorf("test error %w", domain.ErrNotFound)
	require.True(t, errors.Is(err, domain.ErrNotFound))
	require.False(t, errors.Is(err, domain.ErrValid))
	require.ErrorContains(t, err, "not found")
}

func TestErrNotValid(t *testing.T) {
	err := fmt.Errorf("test error %w", domain.ErrValid)
	require.True(t, errors.Is(err, domain.ErrValid))
	require.False(t, errors.Is(err, domain.ErrNotFound))
	require.ErrorContains(t, err, "validation error")
}
