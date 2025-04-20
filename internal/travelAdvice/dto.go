package traveladvice

type TravelAdviceRequest struct {
	Cities    []string `json:"cities" binding:"required,min=1"` // city names or IDs
	StartDate string   `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string   `json:"end_date" binding:"required,datetime=2006-01-02"`
}