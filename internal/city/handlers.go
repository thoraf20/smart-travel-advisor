package city

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

func GetCities(c *gin.Context) {
	var cities []models.City
	if err := db.DB.Order("name ASC").Find(&cities).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch cities")
		return
	}

	var results []gin.H
	for _, city := range cities {
		results = append(results, gin.H{
			"id":          city.ID,
			"name":        city.Name,
			"country":     city.Country,
			"description": city.Description,
			"latitude":    city.Latitude,
			"longitude":   city.Longitude,
		})
	}

	response.Success(c, "Cities retrieved successfully", results)
}
