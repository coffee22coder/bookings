package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coffee22coder/bookings/internal/adapter/http/handler"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	gotError bool
}

func (r *fakeRepo) List(ctx context.Context, limit int, offset int) ([]domain.Airport, error) {
	if r.gotError {
		return nil, errors.New("db down")
	}
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

func TestAirportHandler_InvalidLimit(t *testing.T) {
	repo := &fakeRepo{}
	service := service.AirportServiceNew(repo)
	airportHandler := handler.NewAirportHandler(service)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/airports?limit=-1&offset=0", nil)
	airportHandler.List(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"status":"down", "error": "Invalid limit"}`, rec.Body.String())

}
