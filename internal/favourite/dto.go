package favorite

type AddFavoriteRequest struct {
	TravelAdviceID string `json:"travel_advice_id" binding:"required,uuid"`
}
