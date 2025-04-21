package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/binding"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
	"github.com/thoraf20/smart-travel-advisor/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.RegisterRequest true "Register info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func RegisterHandler(c *gin.Context) {
	var req *RegisterRequest
	req, err := binding.StrictBindJSON[RegisterRequest](c)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	var existing models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; 
	err == nil {
		response.Error(c, http.StatusConflict, "Email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error hashing password", err.Error())
		return
	}

	newUser := models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	if err := db.DB.Create(&newUser).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}
	response.Created(c, "User registered successfully", gin.H{"user_id": newUser.ID})
}

// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.LoginRequest true "Login info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var req *LoginRequest
	req, err := binding.StrictBindJSON[LoginRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	var user *models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); 
	err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.Success(c, "Login successful", gin.H{"token": token})
}

// @Summary Request password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.PasswordResetRequest true "Password reset info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/password/reset/request [post]
func PasswordResetRequestHandler(c *gin.Context) {
	var req *PasswordResetRequest

	req, err := binding.StrictBindJSON[PasswordResetRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	var user *models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	// TODO: Implement password reset logic (e.g., send email with reset link)

	response.Success(c, "Password reset request processed", gin.H{"message": "use this static code for password reset", "reset code": "123456"})
}

// @Summary Reset reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.PasswordResetRequest true "Reset password info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/password/reset [post]
func ResetPasswordRequestHandler(c *gin.Context) {
	var req *ResetPasswordRequest

	req, err := binding.StrictBindJSON[ResetPasswordRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	code := req.Code
	if code != "123456" {
		response.Error(c, http.StatusBadRequest, "Invalid reset code")
		return
	}

	var user *models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; 
	err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error hashing new password", err.Error())
		return
	}

	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := db.DB.Save(&user).Error;
		err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update password", err.Error())
		return
	}

	response.Success(c, "Password updated successfully", gin.H{"message": "Password updated successfully"})
}