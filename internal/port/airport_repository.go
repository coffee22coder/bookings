package port

import (
	"context"

	"github.com/coffee22coder/bookings/internal/domain"
)

type AirportRepository interface {
	List(ctx context.Context, limit int, offset int) ([]domain.Airport, error)
	Count(ctx context.Context) (int, error)
}
