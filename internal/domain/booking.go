package domain

import "time"

type Segment struct {
	FlightID int64
	Fare     string
	Price    float64
}

type Passenger struct {
	TicketNo string
	Name     string
	Document string
	Segments []Segment
}

type Booking struct {
	BookRef     string
	TotalAmount float64
	BookDate    time.Time
	Passengers  []Passenger
}
