package preferences

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
	"gorm.io/gorm/clause"
)

func GetPreferences(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var prefs models.UserPreferences
	err := db.DB.First(&prefs, "user_id = ?", userID).Error
	if err != nil {
		response.Success(c, "No preferences found", gin.H{"preferences": gin.H{}})
		return
	}

	var parsed map[string]interface{}
	_ = json.Unmarshal(prefs.Preferences, &parsed)

	response.Success(c, "Preferences retrieved", gin.H{"preferences": parsed})
}

func UpdatePreferences(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); 
	err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	data, _ := json.Marshal(input)
	pref := models.UserPreferences{
		UserID:     uuid.MustParse(userID),
		Preferences: data,
	}

	// Upsert	
	err := db.DB.
	Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).
	Create(&pref).Error

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save preferences")
		return
	}
	
	response.Success(c, "Preferences updated", gin.H{"preferences": input})
}
