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
	if limit == 0 {
		limit = 20
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *AirportService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
