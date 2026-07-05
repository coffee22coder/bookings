package domain

import "time"

type Flight struct {
	FlightID           int64
	RouteNo            string
	Status             string
	DepartureAirport   string
	ArrivalAirport     string
	ScheduledDeparture time.Time
	ScheduledArrival   time.Time
	ActualDeparture    *time.Time
	ActualArrival      *time.Time
}

type FlightsCount struct {
	Total int
}
