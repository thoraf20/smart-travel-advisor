package user

type UpdateAccountRequest struct {
	//omitempty --> only validate if the field is present
	Name  string `json:"name" binding:"omitempty,min=2"`
	Email string `json:"email" binding:"omitempty,email"`
}