package auth

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateAccountRequest struct {
	//omitempty --> only validate if the field is present
	Name  string `json:"name" binding:"omitempty,min=2"`
	Email string `json:"email" binding:"omitempty,email"`
}