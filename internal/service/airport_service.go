package service

import (
	"context"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/port"
)

type AirportService struct {
	repo port.AirportRepository
}

func New(repo port.AirportRepository) *AirportService {
	return &AirportService{
		repo: repo,
	}
}

func (s *AirportService) List(ctx context.Context, limit int, offset int) ([]domain.Airport, error) {
	return s.repo.List(ctx, limit, offset)
}
