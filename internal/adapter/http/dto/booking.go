package dto

import "time"

type CreateBookingRequest struct {
	TotalAmount float64            `json:"total_amount"`
	Passengers  []PassengerRequest `json:"passengers"`
}

type PassengerRequest struct {
	Name     string           `json:"name"`
	Document string           `json:"document"`
	Segments []SegmentRequest `json:"segments"`
}

type SegmentRequest struct {
	FlightID int64   `json:"flight_id"`
	Fare     string  `json:"fare"`
	Price    float64 `json:"price"`
}

type CreateBookingResponce struct {
	BookRef     string              `json:"bookRef"`
	TotalAmount float64             `json:"totalAmount"`
	BookDate    time.Time           `json:"bookDate"`
	Passengers  []PassengerResponce `json:"passengers"`
}

type PassengerResponce struct {
	TicketNo string            `json:"ticketNo"`
	Name     string            `json:"name"`
	Document string            `json:"document"`
	Segments []SegmentResponce `json:"segments"`
}

type SegmentResponce struct {
	FlightID int64   `json:"flightId"`
	Fare     string  `json:"fare"`
	Price    float64 `json:"price"`
}
