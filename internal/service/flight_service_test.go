package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/coffee22coder/bookings/internal/testutil"
	"github.com/stretchr/testify/require"
)

type fakeFlightRepo struct {
	gotLimit  int
	gotOffset int
	gotID     int64
}

func (r *fakeFlightRepo) List(
	ctx context.Context,
	from string,
	to string,
	date string,
	limit int,
	offset int) ([]domain.Flight, error) {
	r.gotLimit = limit
	r.gotOffset = offset
	d, _ := time.Parse("2006-01-02", testutil.Date)

	fakeAirports := []domain.Flight{
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
	return fakeAirports, nil
}

func (r *fakeFlightRepo) CountSearch(ctx context.Context, from string, to string, date string) (int, error) {
	return 5, nil
}

func (r *fakeFlightRepo) GetByID(ctx context.Context, id int64) (*domain.Flight, error) {
	if id == 777 {
		return nil, domain.ErrNotFound
	}

	now := time.Now()
	r.gotID = id

	fakeFlight := &domain.Flight{
		FlightID:           testutil.FlightID,
		RouteNo:            testutil.RouteNo,
		Status:             testutil.Status,
		DepartureAirport:   testutil.DepartureAirport,
		ArrivalAirport:     testutil.ArrivalAirport,
		ScheduledDeparture: now,
		ScheduledArrival:   now,
		ActualDeparture:    &now,
		ActualArrival:      &now,
	}

	return fakeFlight, nil
}

type testCase struct {
	name      string
	wantError bool
	gotFrom   string
	gotTo     string
	gotDate   string
	gotLimit  string
	gotOffset string
}

func TestFlightService_Validation(t *testing.T) {
	tests := []testCase{
		{

			gotFrom:   "asD",
			gotTo:     "qWE",
			gotDate:   "2025-10-01",
			gotLimit:  "5",
			gotOffset: "1",
		},
		{
			name:      "OK default limit and offset",
			wantError: false,
			gotFrom:   "asD",
			gotTo:     "qWE",
			gotDate:   "2025-10-01",
			gotLimit:  "",
			gotOffset: "",
		},
		{
			name:      "Error: invalid limit int",
			wantError: true,
			gotFrom:   "asD",
			gotTo:     "qWE",
			gotDate:   "2025-10-01",
			gotLimit:  "101",
			gotOffset: "",
		},
		{
			name:      "Error: invalid offset int",
			wantError: true,
			gotFrom:   "asD",
			gotTo:     "qWE",
			gotDate:   "2025-10-01",
			gotLimit:  "100",
			gotOffset: "-1",
		},
		{
			name:      "Error: invalid offset int",
			wantError: true,
			gotFrom:   "",
			gotTo:     "qWE",
			gotDate:   "2025-10-01",
			gotLimit:  "100",
			gotOffset: "12",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &fakeFlightRepo{}
			flightService := service.FlightServiceNew(repo)

			flights, err := flightService.List(ctx, tc.gotFrom, tc.gotTo, tc.gotDate, tc.gotLimit, tc.gotOffset)
			if tc.wantError {
				require.Error(t, err)
				require.True(t, errors.Is(err, domain.ErrValid))
			} else {
				if err != nil {
					require.NoError(t, err)
				}

				if tc.gotLimit == "" && tc.gotOffset == "" {
					require.Equal(t, 20, repo.gotLimit)
					require.Equal(t, 0, repo.gotOffset)
				} else {
					require.Equal(t, 5, repo.gotLimit)
					require.Equal(t, 1, repo.gotOffset)
				}
				require.Len(t, flights, 2)
			}
		})
	}
}

type TestCaseGetById struct {
	name      string
	wantError bool
	gotID     string
}

func TestFlightService_GetById(t *testing.T) {
	tests := []TestCaseGetById{
		{
			name:      "OK valid",
			wantError: false,
			gotID:     "1123",
		},
		{
			name:      "flight ID invalid",
			wantError: true,
			gotID:     "0",
		},
		{
			name:      "flight ID not found",
			wantError: true,
			gotID:     "777",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &fakeFlightRepo{}

			flightService := service.FlightServiceNew(repo)

			flight, err := flightService.GetByID(ctx, tc.gotID)
			if tc.wantError {
				require.Error(t, err)
				if tc.gotID == "777" {
					require.True(t, errors.Is(err, domain.ErrNotFound))
				} else {
					require.True(t, errors.Is(err, domain.ErrValid))
				}

			} else {
				if err != nil {
					require.NoError(t, err)
				}

				intId, err := strconv.ParseInt(tc.gotID, 10, 64)
				if err != nil {
					require.NoError(t, err)
				}

				require.Equal(t, intId, repo.gotID)

				require.Equal(t, testutil.FlightID, flight.FlightID)
			}
		})
	}
}

var errRepoFailed = errors.New("repo failed")

type errorFlightRepo struct{}

func (r *errorFlightRepo) List(
	ctx context.Context,
	from string,
	to string,
	date string,
	limit int,
	offset int) ([]domain.Flight, error) {
	return nil, errRepoFailed
}

func (r *errorFlightRepo) CountSearch(
	ctx context.Context,
	from string,
	to string,
	date string) (int, error) {
	return 0, errRepoFailed
}

func (r *errorFlightRepo) GetByID(
	ctx context.Context,
	id int64) (*domain.Flight, error) {
	return nil, errRepoFailed
}

func TestFlightService_Error(t *testing.T) {
	ctx := context.Background()
	repo := &errorFlightRepo{}
	flightService := service.FlightServiceNew(repo)
	flights, err := flightService.List(ctx, "AER", "FRA", "2025-10-01", "", "")
	require.ErrorIs(t, err, errRepoFailed)
	require.ErrorContains(t, err, "repo failed")
	require.Nil(t, flights)
}
