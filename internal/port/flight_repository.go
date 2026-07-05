package port

import (
	"context"

	"github.com/coffee22coder/bookings/internal/domain"
)

type FlightRepository interface {
	List(
		ctx context.Context,
		from string,
		to string,
		date string,
		limit int,
		offset int) ([]domain.Flight, error)

	CountSearch(
		ctx context.Context,
		from string,
		to string,
		date string) (int, error)
}
