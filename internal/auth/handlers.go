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

func GetAccount(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, "Account retrieved successfully", gin.H{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	})
}

func UpdateAccount(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	req, err := binding.StrictBindJSON[UpdateAccountRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now().Local().String()

	if err := db.DB.Save(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update account")
		return
	}

	response.Success(c, "Account updated successfully", nil)
}

func DeleteAccount(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	if err := db.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	response.Success(c, "Account deleted successfully", nil)
}