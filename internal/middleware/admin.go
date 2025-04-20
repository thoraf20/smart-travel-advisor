package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		adminIDs := os.Getenv("ADMIN_IDS") // e.g. comma-separated UUIDs

		if !containsID(adminIDs, userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func containsID(csv string, userID string) bool {
	// simple match (in real world, store role in DB)
	for _, id := range strings.Split(csv, ",") {
		if id == userID {
			return true
		}
	}
	return false
}
