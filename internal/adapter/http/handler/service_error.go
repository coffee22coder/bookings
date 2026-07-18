package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coffee22coder/bookings/internal/adapter/http/dto"
	"github.com/coffee22coder/bookings/internal/domain"
)

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
