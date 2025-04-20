package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/binding"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

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

func GetAccountQuota(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var count int64
	err := db.DB.Model(&models.TravelAdvice{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve quota usage")
		return
	}

	response.Success(c, "Quota retrieved successfully", gin.H{
		"advice_requests": count,
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
