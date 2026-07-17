package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
	"github.com/go-chi/chi/v5"
)

type FlightHandler struct {
	service *service.FlightService
}

func NewFlightHandler(service *service.FlightService) *FlightHandler {
	return &FlightHandler{
		service: service,
	}
}

func (h *FlightHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var limit, offset, from, to, date string
	limit = r.URL.Query().Get("limit")
	offset = r.URL.Query().Get("offset")
	from = r.URL.Query().Get("from")
	to = r.URL.Query().Get("to")
	date = r.URL.Query().Get("date")

	flights, err := h.service.List(r.Context(), from, to, date, limit, offset)

	if err != nil {
		writeServiceError(w, err)
		return

	}

	flightsRes := make([]dto.FlightResponse, 0, len(flights))

	for _, flight := range flights {
		flightsRes = append(flightsRes, dto.FlightResponse{
			FlightID:           flight.FlightID,
			RouteNo:            flight.RouteNo,
			Status:             flight.Status,
			DepartureAirport:   flight.DepartureAirport,
			ArrivalAirport:     flight.ArrivalAirport,
			ScheduledDeparture: flight.ScheduledDeparture,
			ScheduledArrival:   flight.ScheduledArrival,
			ActualDeparture:    flight.ActualDeparture,
			ActualArrival:      flight.ActualArrival,
		})
	}

	count, err := h.service.CountSearch(r.Context(), from, to, date)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(dto.FlightList{Items: flightsRes, Total: count})
}

func (h *FlightHandler) CountSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var from, to, date string
	from = r.URL.Query().Get("from")
	to = r.URL.Query().Get("to")
	date = r.URL.Query().Get("date")

	count, err := h.service.CountSearch(r.Context(), from, to, date)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(dto.FlightsCount{Total: count})
}

func (h *FlightHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	flightID := chi.URLParam(r, "id")

	flight, err := h.service.GetByID(r.Context(), flightID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	flightRes := dto.FlightResponse{
		FlightID:           flight.FlightID,
		RouteNo:            flight.RouteNo,
		Status:             flight.Status,
		DepartureAirport:   flight.DepartureAirport,
		ArrivalAirport:     flight.ArrivalAirport,
		ScheduledDeparture: flight.ScheduledDeparture,
		ScheduledArrival:   flight.ScheduledArrival,
		ActualDeparture:    flight.ActualDeparture,
		ActualArrival:      flight.ActualArrival,
	}
	json.NewEncoder(w).Encode(flightRes)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValid):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(dto.JSONResponse{
		Status:     "down",
		ErrMessage: err.Error(),
	})
}
