package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
