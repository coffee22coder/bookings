package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
	"github.com/coffee22coder/bookings/internal/service"
)

type BookingHandler struct {
	service *service.BookingSevice
}

func NewBookingHandler(service *service.BookingSevice) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeServiceError(w, fmt.Errorf("invalid json: %w", domain.ErrValid))
		return
	}

	b, err := h.service.Create(r.Context(), req.TotalAmount, req.Passengers)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	passengers := make([]dto.PassengerResponce, 0, len(b.Passengers))

	for _, p := range b.Passengers {
		segments := make([]dto.SegmentResponce, 0, len(p.Segments))
		for _, s := range p.Segments {
			ns := dto.SegmentResponce{
				FlightID: s.FlightID,
				Price:    s.Price,
				Fare:     s.Fare,
			}
			segments = append(segments, ns)
		}

		np := dto.PassengerResponce{
			TicketNo: p.TicketNo,
			Name:     p.Name,
			Document: p.Document,
			Segments: segments,
		}

		passengers = append(passengers, np)
	}

	bookingRefRes := dto.CreateBookingResponce{
		BookRef:     b.BookRef,
		TotalAmount: b.TotalAmount,
		BookDate:    b.BookDate,
		Passengers:  passengers,
	}

	json.NewEncoder(w).Encode(bookingRefRes)

}
