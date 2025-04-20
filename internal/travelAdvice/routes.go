package traveladvice

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func TravelAdviceRoutes(r *gin.Engine) {
	advice := r.Group("/api/v1/travel-advice")
	advice.Use(middleware.AuthMiddleware())
	{
		advice.POST("", CreateTravelAdvice)
		advice.GET("/history", GetTravelAdviceHistory)
		advice.GET("/history/:id", GetTravelAdviceHistoryById)
		advice.GET("/search", SearchTravelAdviceHistory)
	}
}
