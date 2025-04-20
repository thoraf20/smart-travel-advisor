package auth

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordResetRequest struct {
	Email    string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	NewPassword 	 string `json:"new_password" binding:"required,min=8"`
	Code 	 string `json:"code" binding:"required,min=6"`
}
