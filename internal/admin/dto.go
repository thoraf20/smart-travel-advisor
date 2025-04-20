package admin

type CreateCityRequest struct {
	Name        string  `json:"name" binding:"required"`
	Country     string  `json:"country" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type UpdateCityRequest struct {
	Name        *string  `json:"name,omitempty"`
	Country     *string  `json:"country,omitempty"`
	Description *string  `json:"description,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}
