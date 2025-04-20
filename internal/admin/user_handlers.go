package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Order("created_at DESC").Find(&users).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	response.Success(c, "Users retrieved", users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, "id = ?", id).Error; 
	err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}
	response.Success(c, "User retrieved", user)
}

func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); 
	err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", id).Updates(input).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.Success(c, "User updated", nil)
}

func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.User{}, "id = ?", id).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	response.Success(c, "User deleted", nil)
}
