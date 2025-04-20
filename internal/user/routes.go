package user

import (
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/smart-travel-advisor/internal/middleware"
)

func UserRoutes(r *gin.Engine) {

	account := r.Group("api/v1/user/account")
	account.Use(middleware.AuthMiddleware())
	{
		account.GET("", GetAccountQuota)
		account.GET("/quota", GetAccountQuota) // ✅ Added
		account.PUT("", UpdateAccount)
		account.DELETE("", DeleteAccount)
	}
}
