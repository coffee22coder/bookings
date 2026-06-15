package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httpadapter "github.com/coffee22coder/bookings/internal/adapter/http"
	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/config"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/stretchr/testify/require"
)

type fakeDB struct {
	err error
}

func (f fakeDB) Ping(ctx context.Context) error {
	return f.err
}

type testCaseHealth struct {
	name           string
	method         string
	path           string
	reqID          string
	wantStatusCode int
	wantJSON       string
	wantReqID      string
}

func TestServerHealth(t *testing.T) {
	tests := []testCaseHealth{
		{
			name:           "Returns200",
			method:         http.MethodGet,
			path:           "/health",
			wantStatusCode: http.StatusOK,
			wantJSON:       `{"status":"ok"}`,
		},
		{
			name:           "Returns405",
			method:         http.MethodPost,
			path:           "/health",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "ReturnsJSON",
			method:         http.MethodGet,
			path:           "/health",
			wantStatusCode: http.StatusOK,
			wantJSON:       `{"status":"ok"}`,
		},
		{
			name:           "RequestID",
			method:         http.MethodGet,
			path:           "/health",
			wantStatusCode: http.StatusOK,
			wantJSON:       `{"status":"ok"}`,
			reqID:          "test-123",
			wantReqID:      "test-123",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{HTTPPort: 8080}
			srv := httpadapter.NewServer(cfg, logger, fakeDB{}, &httpadapter.Services{})
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.reqID != "" {
				req.Header.Set("X-Request-Id", tc.reqID)
			}

			srv.Handler().ServeHTTP(rec, req)

			require.Equal(t, tc.wantStatusCode, rec.Code)

			if tc.wantJSON != "" {
				require.JSONEq(t, tc.wantJSON, rec.Body.String())
			}

			if tc.wantReqID != "" {
				require.Equal(t, tc.wantReqID, rec.Header().Get("X-Request-Id"))
			}
		})
	}
}

type testCaseReady struct {
	name           string
	dbErr          error
	wantStatusCode int
	wantJSON       string
}

func TestServerReady(t *testing.T) {
	tests := []testCaseReady{
		{
			name:           "DB ping Error",
			dbErr:          errors.New("down"),
			wantStatusCode: http.StatusServiceUnavailable,
			wantJSON:       `{"status":"down", "error": "down"}`,
		},
		{
			name:           "DB ping OK",
			dbErr:          nil,
			wantStatusCode: http.StatusOK,
			wantJSON:       `{"status":"ok"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{HTTPPort: 8080}
			var db fakeDB
			if tc.dbErr != nil {
				db = fakeDB{err: errors.New("down")}
			} else {
				db = fakeDB{}
			}
			srv := httpadapter.NewServer(cfg, logger, db, &httpadapter.Services{})
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ready", nil)

			srv.Handler().ServeHTTP(rec, req)

			require.Equal(t, tc.wantStatusCode, rec.Code)

			if tc.wantJSON != "" {
				require.JSONEq(t, tc.wantJSON, rec.Body.String())
			}
		})
	}
}

func TestRecoverMiddleware(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{HTTPPort: 8080}
	srv := httpadapter.NewServer(cfg, logger, fakeDB{}, &httpadapter.Services{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	srv.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

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

func TestServerAirports_OK(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{HTTPPort: 8080}
	var db fakeDB
	srv := httpadapter.NewServer(cfg, logger, db, &httpadapter.Services{
		Airports: service.New(&fakeRepo{gotError: false}),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/airports", nil)

	srv.Handler().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var got dto.AerportList
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	require.NoError(t, err)
	require.Len(t, got.Items, 3)
}

func TestServerAirports_ErrorPath(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{HTTPPort: 8080}
	var db fakeDB
	srv := httpadapter.NewServer(cfg, logger, db, &httpadapter.Services{
		Airports: service.New(&fakeRepo{gotError: true}),
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/airports", nil)

	srv.Handler().ServeHTTP(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.JSONEq(t, `{"status":"down", "error": "db down"}`, rec.Body.String())
}
