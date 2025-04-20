package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func RegisterAdminRoutes(r *gin.Engine) {
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/users", GetAllUsers)
		admin.GET("/users/:id", GetUserByID)
		admin.PUT("/users/:id", UpdateUserByID)
		admin.DELETE("/users/:id", DeleteUserByID)
		admin.POST("/cities", CreateCity)
		admin.PUT("/cities/:id", UpdateCity)
		admin.DELETE("/cities/:id", DeleteCity)
	}
}
