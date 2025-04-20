package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("api/v1/auth")
	auth.POST("/signup", RegisterHandler)
	auth.POST("/login", LoginHandler)
	auth.POST("/password/reset/request", ResetPasswordRequestHandler)
	auth.POST("/password/reset", PasswordResetRequestHandler)
}
