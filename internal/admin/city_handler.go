package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/binding"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

func CreateCity(c *gin.Context) {
	req, err := binding.StrictBindJSON[CreateCityRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	city := models.City{
		ID:          uuid.New(),
		Name:        req.Name,
		Country:     req.Country,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := db.DB.Create(&city).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create city")
		return
	}

	response.Created(c, "City created", gin.H{"city_id": city.ID})
}

func UpdateCity(c *gin.Context) {
	id := c.Param("id")
	var city models.City
	if err := db.DB.First(&city, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "City not found")
		return
	}

	req, err := binding.StrictBindJSON[UpdateCityRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	if req.Name != nil {
		city.Name = *req.Name
	}
	if req.Country != nil {
		city.Country = *req.Country
	}
	if req.Description != nil {
		city.Description = *req.Description
	}
	if req.Latitude != nil {
		city.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		city.Longitude = *req.Longitude
	}

	if err := db.DB.Save(&city).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update city")
		return
	}

	response.Success(c, "City updated", nil)
}

func DeleteCity(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.City{}, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete city")
		return
	}
	response.Success(c, "City deleted", nil)
}
