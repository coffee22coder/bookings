package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/port"
	"github.com/coffee22coder/bookings/internal/util"
)

func validFlightFilter(from, to, date string) error {
	if from == "" || to == "" || date == "" {
		return fmt.Errorf("invalid params (from, to, date): %w", domain.ErrValid)
	}

	from = strings.ToUpper(strings.TrimSpace(from))
	to = strings.ToUpper(strings.TrimSpace(to))

	if !util.IsIATA(from) || !util.IsIATA(to) {
		return fmt.Errorf("invalid airprot IATA code (from, to): %w", domain.ErrValid)
	}

	date = strings.TrimSpace(date)
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date (layout 2006-01-02): %w", domain.ErrValid)
	}

	return nil
}

type FlightService struct {
	repo port.FlightRepository
}

func FlightServiceNew(repo port.FlightRepository) *FlightService {
	return &FlightService{
		repo: repo,
	}
}

func (s *FlightService) List(
	ctx context.Context,
	from string,
	to string,
	date string,
	limit string,
	offset string) ([]domain.Flight, error) {
	var limitInt, offsetInt int
	var err error
	if limit == "" {
		limitInt = 20
	} else {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("invalid limit: %w", domain.ErrValid)
		}
		if limitInt < 0 || limitInt > 100 {
			return nil, fmt.Errorf("invalid limit: %w", domain.ErrValid)
		}
	}
	if offset == "" {
		offsetInt = 0
	} else {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			return nil, fmt.Errorf("invalid offset: %w", domain.ErrValid)
		}
		if offsetInt < 0 {
			return nil, fmt.Errorf("invalid offsetInt: %w", domain.ErrValid)
		}
	}

	if from == "" || to == "" || date == "" {
		return nil, fmt.Errorf("invalid params (from, to, date): %w", domain.ErrValid)
	}

	err = validFlightFilter(from, to, date)
	if err != nil {
		return nil, err
	}

	return s.repo.List(
		ctx,
		from,
		to,
		date,
		limitInt,
		offsetInt)
}

func (s *FlightService) CountSearch(
	ctx context.Context,
	from string,
	to string,
	date string) (int, error) {
	err := validFlightFilter(from, to, date)
	if err != nil {
		return 0, err
	}

	return s.repo.CountSearch(
		ctx,
		from,
		to,
		date)
}
