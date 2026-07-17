package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coffee22coder/bookings/internal/adapter/http/handler"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/coffee22coder/bookings/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type fakeFlightRepo struct {
	gotError bool
}

func (r *fakeFlightRepo) List(
	ctx context.Context,
	from string,
	to string,
	date string,
	limit int,
	offset int) ([]domain.Flight, error) {
	if r.gotError {
		return nil, errors.New("db down")
	}
	d, _ := time.Parse("2006-01-02", testutil.Date)

	fakeFlights := []domain.Flight{
		{
			FlightID:           testutil.FlightID,
			RouteNo:            testutil.RouteNo,
			Status:             testutil.Status,
			DepartureAirport:   testutil.DepartureAirport,
			ArrivalAirport:     testutil.ArrivalAirport,
			ScheduledDeparture: d,
			ScheduledArrival:   d,
			ActualDeparture:    &d,
			ActualArrival:      &d,
		},
		{
			FlightID:           testutil.FlightID,
			RouteNo:            testutil.RouteNo,
			Status:             testutil.Status,
			DepartureAirport:   testutil.DepartureAirport,
			ArrivalAirport:     testutil.ArrivalAirport,
			ScheduledDeparture: d,
			ScheduledArrival:   d,
			ActualDeparture:    &d,
			ActualArrival:      &d,
		},
	}
	return fakeFlights, nil
}

func (r *fakeFlightRepo) CountSearch(ctx context.Context, from string, to string, date string) (int, error) {
	return 5, nil
}

func (r *fakeFlightRepo) GetByID(
	ctx context.Context,
	id int64) (*domain.Flight, error) {
	d, _ := time.Parse("2006-01-02", testutil.Date)
	fakeFlight := &domain.Flight{
		FlightID:           testutil.FlightID,
		RouteNo:            testutil.RouteNo,
		Status:             testutil.Status,
		DepartureAirport:   testutil.DepartureAirport,
		ArrivalAirport:     testutil.ArrivalAirport,
		ScheduledDeparture: d,
		ScheduledArrival:   d,
		ActualDeparture:    &d,
		ActualArrival:      &d,
	}
	return fakeFlight, nil
}

func TestFlightHandler_OK(t *testing.T) {
	repo := &fakeFlightRepo{}
	service := service.FlightServiceNew(repo)
	handler := handler.NewFlightHandler(service)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/flights?from=LED&to=SVO&date=2025-10-01&limit=1&offset=0", nil)
	handler.List(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	router := chi.NewRouter()
	router.Get("/api/v1/flights/{id}", handler.GetByID)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/flights/1123", nil)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

}

func TestFlightHandler_BadRequest(t *testing.T) {
	repo := &fakeFlightRepo{}
	service := service.FlightServiceNew(repo)
	handler := handler.NewFlightHandler(service)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/flights?from=LED&to=SVO&date=2025-10-01&limit=-1&offset=0", nil)
	handler.List(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"status":"down", "error": "invalid limit int: validation error"}`, rec.Body.String())

	router := chi.NewRouter()
	router.Get("/api/v1/flights/{id}", handler.GetByID)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/flights/0", nil)
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"status":"down", "error": "invalid flightId: validation error"}`, rec.Body.String())

}
