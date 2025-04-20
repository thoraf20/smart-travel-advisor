package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thoraf20/smart-travel-advisor/config"
	"github.com/thoraf20/smart-travel-advisor/internal/admin"
	"github.com/thoraf20/smart-travel-advisor/internal/auth"
	"github.com/thoraf20/smart-travel-advisor/internal/city"
	"github.com/thoraf20/smart-travel-advisor/internal/db"
	favorite "github.com/thoraf20/smart-travel-advisor/internal/favourite"
	"github.com/thoraf20/smart-travel-advisor/internal/preferences"
	traveladvice "github.com/thoraf20/smart-travel-advisor/internal/travelAdvice"
	"github.com/thoraf20/smart-travel-advisor/internal/user"
)

func main() {

	config.LoadConfig()

	db.InitDB()

	r := gin.Default()

	auth.RegisterAuthRoutes(r)
	user.UserRoutes(r)
	traveladvice.TravelAdviceRoutes(r)
	favorite.FavoriteRoutes(r)
	city.RegisterCityRoutes(r)
	preferences.RegisterPreferencesRoutes(r)
	admin.RegisterAdminRoutes(r)

	log.Println("🚀 Starting server on port", viper.GetString("PORT"))
	r.Run(":" + viper.GetString("PORT"))
}
