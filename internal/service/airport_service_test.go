package service_test

import (
	"context"
	"testing"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	gotLimit int
}

func (r *fakeRepo) List(ctx context.Context, limit int, offset int) ([]domain.Airport, error) {
	r.gotLimit = limit
	fakeAirports := []domain.Airport{
		{
			Code:        "123",
			Name:        domain.Location{En: "En", Ru: "Ru"},
			City:        domain.Location{En: "En", Ru: "Ru"},
			Country:     domain.Location{En: "En", Ru: "Ru"},
			Coordinates: [2]float64{1.212, 2.212},
			Timezone:    "TZN",
		},
		{
			Code:        "123",
			Name:        domain.Location{En: "En", Ru: "Ru"},
			City:        domain.Location{En: "En", Ru: "Ru"},
			Country:     domain.Location{En: "En", Ru: "Ru"},
			Coordinates: [2]float64{1.212, 2.212},
			Timezone:    "TZN",
		},
		{
			Code:        "123",
			Name:        domain.Location{En: "En", Ru: "Ru"},
			City:        domain.Location{En: "En", Ru: "Ru"},
			Country:     domain.Location{En: "En", Ru: "Ru"},
			Coordinates: [2]float64{1.212, 2.212},
			Timezone:    "TZN",
		},
	}
	return fakeAirports, nil
}

func (r *fakeRepo) Count(ctx context.Context) (int, error) {
	return 5, nil
}

func TestAirportService_List(t *testing.T) {
	ctx := context.Background()

	repo := &fakeRepo{}
	airportService := service.New(repo)

	airports, err := airportService.List(ctx, 3, 0)
	require.NoError(t, err)
	require.Len(t, airports, 3)
}

func TestAirportService_Count(t *testing.T) {
	ctx := context.Background()

	repo := &fakeRepo{}
	airportService := service.New(repo)

	count, err := airportService.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, count, 5)
}

func TestAirportService_DefaultLimit(t *testing.T) {
	ctx := context.Background()

	repo := &fakeRepo{}
	airportService := service.New(repo)

	airports, err := airportService.List(ctx, 0, 0)
	require.NoError(t, err)
	require.Len(t, airports, 3)
	require.Equal(t, 20, repo.gotLimit)
}
