package city

import "github.com/gin-gonic/gin"

func RegisterCityRoutes(r *gin.Engine) {
	r.GET("/api/v1/cities", GetCities)
}
