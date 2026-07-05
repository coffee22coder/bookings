package dto

import "time"

type FlightResponse struct {
	FlightID           int64      `json:"flightID"`
	RouteNo            string     `json:"routeNo"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departureAirport"`
	ArrivalAirport     string     `json:"arrivalAirport"`
	ScheduledDeparture time.Time  `json:"scheduledDeparture"`
	ScheduledArrival   time.Time  `json:"scheduledArrival"`
	ActualDeparture    *time.Time `json:"actualDeparture"`
	ActualArrival      *time.Time `json:"actualArrival"`
}

type FlightList struct {
	Items []FlightResponse `json:"items"`
	Total int              `json:"total"`
}

type FlightsCount struct {
	Total int `json:"total"`
}
