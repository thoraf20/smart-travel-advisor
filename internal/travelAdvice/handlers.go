package traveladvice

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	"github.com/thoraf20/smart-travel-advisor/internal/integrations"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/thoraf20/smart-travel-advisor/pkg/binding"
	"github.com/thoraf20/smart-travel-advisor/pkg/response"
)

func CreateTravelAdvice(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	req, err := binding.StrictBindJSON[TravelAdviceRequest](c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	cityJSON, _ := json.Marshal(req.Cities)

	weather, err := integrations.GetWeather(req.Cities[0])

	if err != nil {
		weather = &integrations.WeatherResponse{
			Weather: []struct {
				Main        string "json:\"main\""
				Description string "json:\"description\""
			}{{Main: "Unavailable", Description: "Could not fetch weather"}},
		}
	}

	flights, err := integrations.GetFlightsArrivingInCity(req.Cities[0]) // simulate with first city
	var flightData []map[string]string

	if err == nil {
		for _, f := range flights {
			flightData = append(flightData, map[string]string{
				"airline":     f.Airline.Name,
				"flight_no":   f.Flight.Number,
				"from":        f.Departure.Airport,
				"to":          f.Arrival.Airport,
				"departure":   f.Departure.Time,
				"arrival":     f.Arrival.Time,
			})
		}	
	} else {
		flightData = []map[string]string{
			{"error": "Unable to fetch flights"},
		}
	}

	adviceResponse := map[string]interface{}{
		"weather":  map[string]interface{}{
			"flights":  weather.Weather[0].Main,
			"description":   weather.Weather[0].Description,
			"temp": weather.Main.Temp,
			"humidity": weather.Main.Humidity,
		},
		"flights":  flightData,
		"tips":     "Dress accordingly 😎",
		"itinerary": []string{"Explore local attractions", "Sample street food"},
	}
	
	adviceJSON, _ := json.Marshal(adviceResponse)

	advice := models.TravelAdvice{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		StartDate: startDate,
		EndDate:   endDate,
		Cities:    cityJSON,
		Advice:    adviceJSON,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&advice).Error; 
	err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save travel advice")
		return
	}

	response.Created(c, "Travel advice created", gin.H{
		"advice_id": advice.ID,
		"advice":    adviceResponse,
	})
}

func GetTravelAdviceHistory(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var history []models.TravelAdvice
	err := db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&history).Error

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve travel advice history")
		return
	}

	var results []gin.H
	for _, h := range history {
		var cities []string
		var adviceData map[string]interface{}

		_ = json.Unmarshal(h.Cities, &cities)
		_ = json.Unmarshal(h.Advice, &adviceData)

		results = append(results, gin.H{
			"id":         h.ID,
			"cities":     cities,
			"start_date": h.StartDate.Format("2006-01-02"),
			"end_date":   h.EndDate.Format("2006-01-02"),
			"advice":     adviceData,
			"created_at": h.CreatedAt,
		})
	}

	response.Success(c, "Travel advice history retrieved", results)
}

func GetTravelAdviceHistoryById(c *gin.Context) {
	adviceID := c.Param("advice_id")
	userID := c.MustGet("userID").(string)

	var history models.TravelAdvice
	err := db.DB.
		Where("id = ? AND user_id = ?", adviceID, userID).
		First(&history).Error

	if err != nil {
		response.Error(c, http.StatusNotFound, "Travel advice not found")
		return
	}

	var cities []string
	var adviceData map[string]interface{}

	_ = json.Unmarshal(history.Cities, &cities)
	_ = json.Unmarshal(history.Advice, &adviceData)

	response.Success(c, "Travel advice retrieved", gin.H{
		"id":         history.ID,
		"cities":     cities,
		"start_date": history.StartDate.Format("2006-01-02"),
		"end_date":   history.EndDate.Format("2006-01-02"),
		"advice":     adviceData,
	})
}

func SearchTravelAdviceHistory(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	page := c.DefaultQuery("page", "1")
	query := c.DefaultQuery("query", "")

	pageNum, _ := strconv.Atoi(page)
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize := 10
	offset := (pageNum - 1) * pageSize

	var records []models.TravelAdvice
	err := db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to search travel advice")
		return
	}

	var filtered []gin.H
	for _, r := range records {
		var cities []string
		var advice map[string]interface{}
		_ = json.Unmarshal(r.Cities, &cities)
		_ = json.Unmarshal(r.Advice, &advice)

		// Filter by query (case-insensitive)
		match := false
		for _, city := range cities {
			if strings.Contains(strings.ToLower(city), strings.ToLower(query)) {
				match = true
				break
			}
		}
		if query == "" || match {
			filtered = append(filtered, gin.H{
				"id":         r.ID,
				"cities":     cities,
				"start_date": r.StartDate.Format("2006-01-02"),
				"end_date":   r.EndDate.Format("2006-01-02"),
				"advice":     advice,
				"created_at": r.CreatedAt,
			})
		}
	}

	response.Success(c, "Filtered travel advice", gin.H{
		"page":     pageNum,
		"pageSize": pageSize,
		"results":  filtered,
	})
}

func GetTravelAdviceByID(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	adviceID := c.Param("id")

	var record models.TravelAdvice
	err := db.DB.Where("id = ? AND user_id = ?", adviceID, userID).First(&record).Error
	if err != nil {
		response.Error(c, http.StatusNotFound, "Travel advice not found")
		return
	}

	var cities []string
	var adviceData map[string]interface{}
	_ = json.Unmarshal(record.Cities, &cities)
	_ = json.Unmarshal(record.Advice, &adviceData)

	response.Success(c, "Travel advice retrieved", gin.H{
		"id":         record.ID,
		"cities":     cities,
		"start_date": record.StartDate.Format("2006-01-02"),
		"end_date":   record.EndDate.Format("2006-01-02"),
		"advice":     adviceData,
		"created_at": record.CreatedAt,
	})
}

