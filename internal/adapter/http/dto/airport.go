package dto

type AirportResponse struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

type AerportList struct {
	Items []AirportResponse `json:"items"`
	Total int               `json:"total"`
}

type AeroportCount struct {
	Count int `json:"count"`
}
