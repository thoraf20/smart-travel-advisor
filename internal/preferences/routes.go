package preferences

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func RegisterPreferencesRoutes(r *gin.Engine) {
	group := r.Group("/preferences")
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("", GetPreferences)
		group.PUT("", UpdatePreferences)
	}
}
