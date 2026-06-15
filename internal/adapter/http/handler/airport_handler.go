package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/service"
)

type AirportHandler struct {
	service *service.AirportService
}

func NewAirportHandler(service *service.AirportService) *AirportHandler {
	return &AirportHandler{
		service: service,
	}
}

func (h *AirportHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit := 5
	offset := 0
	var err error
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
		if limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.JSONResponse{
				Status:     "down",
				ErrMessage: "Invalid limit",
			})
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.JSONResponse{
				Status:     "down",
				ErrMessage: err.Error(),
			})
			return
		}
	}

	airports, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(dto.JSONResponse{
			Status:     "down",
			ErrMessage: err.Error(),
		})
		return
	}

	airportsRes := make([]dto.AirportResponse, 0, len(airports))

	for _, airport := range airports {
		airportsRes = append(airportsRes, dto.AirportResponse{
			Code:     airport.Code,
			Name:     airport.Name.En,
			City:     airport.City.En,
			Timezone: airport.Timezone,
		})
	}

	json.NewEncoder(w).Encode(dto.AerportList{Items: airportsRes, Total: len(airportsRes)})
}

func (h *AirportHandler) Count(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count, err := h.service.Count(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(dto.JSONResponse{
			Status:     "down",
			ErrMessage: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(dto.AeroportCount{Count: count})
}
