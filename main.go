package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thoraf20/smart-travel-advisor/config"
	"github.com/thoraf20/smart-travel-advisor/pkg/firebase"
)

func main() {
	config.LoadConfig()

	firebase.InitFirebase()

	r := gin.Default()

	// TODO: Setup routes

	log.Println("🚀 Starting server on port", viper.GetString("PORT"))
	r.Run(":" + viper.GetString("PORT"))
}
