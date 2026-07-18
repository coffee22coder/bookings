package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/port"
)

type BookingRepo struct {
	db    *Pool
	genID port.GeneratorID
}

func NewBookingRepo(db *Pool, genID port.GeneratorID) *BookingRepo {
	return &BookingRepo{
		db:    db,
		genID: genID,
	}
}

func (r *BookingRepo) Create(
	ctx context.Context,
	totalAmount float64,
	passengers []dto.PassengerRequest) (*domain.Booking, error) {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	bookRef, err := r.genID.BookRef()

	if err != nil {
		return nil, err
	}

	now := time.Now()

	_, err = tx.Exec(ctx,
		`INSERT INTO bookings.bookings (book_ref, book_date, total_amount) VALUES ($1, $2, $3)`,
		bookRef, now, totalAmount,
	)

	if err != nil {
		return nil, err
	}

	bookingPassengers := make([]domain.Passenger, 0, len(passengers))

	for _, p := range passengers {
		var maxTicket int64
		err = tx.QueryRow(ctx,
			`SELECT COALESCE(MAX(ticket_no::bigint), 0) FROM bookings.tickets`,
		).Scan(&maxTicket)
		if err != nil {
			return nil, err
		}
		ticketNo := fmt.Sprintf("%013d", maxTicket+1)

		_, err = tx.Exec(ctx,
			`INSERT INTO bookings.tickets (ticket_no, book_ref, passenger_name, passenger_id, outbound) VALUES ($1, $2, $3, $4, $5)`,
			ticketNo, bookRef, p.Name, p.Document, true)
		if err != nil {
			return nil, err
		}

		seg := make([]domain.Segment, 0, len(p.Segments))

		for _, s := range p.Segments {
			ns := domain.Segment{
				FlightID: s.FlightID,
				Fare:     s.Fare,
				Price:    s.Price,
			}

			seg = append(seg, ns)
			_, err = tx.Exec(ctx,
				`INSERT INTO bookings.segments (ticket_no, flight_id, fare_conditions, price) VALUES ($1, $2, $3, $4)`,
				ticketNo, s.FlightID, s.Fare, s.Price,
			)
			if err != nil {
				return nil, err
			}
		}

		np := domain.Passenger{
			TicketNo: ticketNo,
			Name:     p.Name,
			Document: p.Document,
			Segments: seg,
		}

		bookingPassengers = append(bookingPassengers, np)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	booking := domain.Booking{
		TotalAmount: totalAmount,
		Passengers:  bookingPassengers,
		BookDate:    now,
		BookRef:     bookRef,
	}

	return &booking, nil

}
