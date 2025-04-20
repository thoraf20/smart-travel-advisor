package favorite

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/binding"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

func AddFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	req, err := binding.StrictBindJSON[AddFavoriteRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Optional: verify travel advice exists
	var advice models.TravelAdvice
	if err := db.DB.First(&advice, "id = ? AND user_id = ?", req.TravelAdviceID, userID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Travel advice not found")
		return
	}

	favorite := models.Favorite{
		ID:             uuid.New(),
		UserID:         uuid.MustParse(userID),
		TravelAdviceID: advice.ID,
		CreatedAt:      time.Now(),
	}

	if err := db.DB.Create(&favorite).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save favorite")
		return
	}

	response.Created(c, "Favorite saved", gin.H{"favorite_id": favorite.ID})
}

func GetFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var favorites []models.Favorite
	if err := db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve favorites")
		return
	}

	var results []gin.H
	for _, fav := range favorites {
		var advice models.TravelAdvice
		if err := db.DB.First(&advice, "id = ?", fav.TravelAdviceID).Error; err != nil {
			continue // skip invalid refs
		}

		var cities []string
		var adviceData map[string]interface{}
		_ = json.Unmarshal(advice.Cities, &cities)
		_ = json.Unmarshal(advice.Advice, &adviceData)

		results = append(results, gin.H{
			"id":          fav.ID,
			"advice_id":   advice.ID,
			"cities":      cities,
			"start_date":  advice.StartDate.Format("2006-01-02"),
			"end_date":    advice.EndDate.Format("2006-01-02"),
			"advice":      adviceData,
			"favorited_at": fav.CreatedAt,
		})
	}

	response.Success(c, "Favorites retrieved", results)
}

func DeleteFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	favoriteID := c.Param("id")

	if err := db.DB.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&models.Favorite{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete favorite")
		return
	}

	response.Success(c, "Favorite deleted", nil)
}
