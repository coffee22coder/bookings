package service

import (
	"context"
	"fmt"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/port"
)

type BookingSevice struct {
	repo port.BookingRepository
}

func BookingServiceNew(repo port.BookingRepository) *BookingSevice {
	return &BookingSevice{
		repo: repo,
	}
}

func (s *BookingSevice) Create(
	ctx context.Context,
	totalAmount float64,
	passengers []dto.PassengerRequest) (*domain.Booking, error) {

	if totalAmount <= 0 {
		return nil, fmt.Errorf("invalid totalAmount: %w", domain.ErrValid)
	}

	if len(passengers) == 0 {
		return nil, fmt.Errorf("invalid passengers: %w", domain.ErrValid)
	}

	return s.repo.Create(ctx, totalAmount, passengers)
}
