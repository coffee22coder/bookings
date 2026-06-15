package dto

type JSONResponse struct {
	Status     string `json:"status"`
	ErrMessage string `json:"error,omitempty"`
}
