package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/jackc/pgx/v5"
)

type FlightRepo struct {
	db *Pool
}

func NewFlightRepo(db *Pool) *FlightRepo {
	return &FlightRepo{
		db: db,
	}
}

func (r *FlightRepo) List(
	ctx context.Context,
	from string,
	to string,
	date string,
	limit int,
	offset int) ([]domain.Flight, error) {
	list := make([]domain.Flight, 0, limit)
	rows, err := r.db.pool.Query(ctx, `SELECT
		f.flight_id,
		f.route_no,
		f.status,
		r.departure_airport,
		r.arrival_airport,
		f.scheduled_departure,
		f.scheduled_arrival,
		f.actual_departure,
		f.actual_arrival 
	FROM bookings.flights f 
	JOIN bookings.routes r ON r.route_no = f.route_no
	WHERE r.departure_airport = $3 AND r.arrival_airport = $4
	AND f.scheduled_departure >= $5::date
	AND f.scheduled_departure <  $5::date + interval '1 day'
	ORDER BY f.scheduled_departure
	LIMIT $1 OFFSET $2;`,
		limit, offset, from, to, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var flightID int64
		var routeNo, status, departureAirport, arrivalAirport string
		var scheduledDeparture, scheduledArrival time.Time
		var actualDeparture, actualArrival *time.Time

		if err := rows.Scan(
			&flightID,
			&routeNo,
			&status,
			&departureAirport,
			&arrivalAirport,
			&scheduledDeparture,
			&scheduledArrival,
			&actualDeparture,
			&actualArrival,
		); err != nil {
			return nil, fmt.Errorf("scan flights: %w", err)
		}
		flight := domain.Flight{
			FlightID:           flightID,
			RouteNo:            routeNo,
			Status:             status,
			DepartureAirport:   departureAirport,
			ArrivalAirport:     arrivalAirport,
			ScheduledDeparture: scheduledDeparture,
			ScheduledArrival:   scheduledArrival,
			ActualDeparture:    actualDeparture,
			ActualArrival:      actualArrival,
		}
		list = append(list, flight)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *FlightRepo) CountSearch(
	ctx context.Context,
	from string,
	to string,
	date string,
) (int, error) {
	var total int
	err := r.db.pool.QueryRow(
		ctx, `SELECT COUNT(DISTINCT(f.flight_id))
	FROM bookings.flights f 
	JOIN bookings.routes r ON r.route_no = f.route_no
	WHERE r.departure_airport = $1 AND r.arrival_airport = $2
	AND f.scheduled_departure >= $3::date
	AND f.scheduled_departure <  $3::date + interval '1 day';`, from, to, date).Scan(&total)
	return total, err
}

func (r *FlightRepo) GetByID(
	ctx context.Context,
	id int64,
) (*domain.Flight, error) {
	var routeNo, status, departureAirport, arrivalAirport string
	var scheduledDeparture, scheduledArrival time.Time
	var actualDeparture, actualArrival *time.Time

	err := r.db.pool.QueryRow(ctx, `SELECT 
		f.route_no,
		f.status,
		r.departure_airport,
		r.arrival_airport,
		f.scheduled_departure,
		f.scheduled_arrival,
		f.actual_departure,
		f.actual_arrival FROM bookings.flights f
		JOIN bookings.routes r ON f.route_no = r.route_no
		WHERE f.flight_id = $1`, id).Scan(
		&routeNo,
		&status,
		&departureAirport,
		&arrivalAirport,
		&scheduledDeparture,
		&scheduledArrival,
		&actualDeparture,
		&actualArrival,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	flight := domain.Flight{
		FlightID:           id,
		RouteNo:            routeNo,
		Status:             status,
		DepartureAirport:   departureAirport,
		ArrivalAirport:     arrivalAirport,
		ScheduledDeparture: scheduledDeparture,
		ScheduledArrival:   scheduledArrival,
		ActualDeparture:    actualDeparture,
		ActualArrival:      actualArrival,
	}
	return &flight, nil
}
