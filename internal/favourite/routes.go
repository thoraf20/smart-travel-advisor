package favorite

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func FavoriteRoutes(r *gin.Engine) {
	fav := r.Group("/favorites")
	fav.Use(middleware.AuthMiddleware())
	{
		fav.POST("", AddFavorite)
		fav.GET("", GetFavorites)
		fav.DELETE("/:id", DeleteFavorite)
	}
}
