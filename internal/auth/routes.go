package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("api/v1/auth")
	auth.POST("/signup", RegisterHandler)
	auth.POST("/login", LoginHandler)

	account := r.Group("api/v1/account")
	account.Use(middleware.AuthMiddleware())
	{
		account.GET("", GetAccount)
		account.PUT("", UpdateAccount)
		account.DELETE("", DeleteAccount)
	}
}
