package port

import (
	"context"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
)

type BookingRepository interface {
	Create(
		ctx context.Context,
		totalAmount float64,
		passengers []dto.PassengerRequest) (*domain.Booking, error)
}
