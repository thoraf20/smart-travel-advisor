package db

import (
	"log"

	"github.com/spf13/viper"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := viper.GetString("DB_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	_ = DB.AutoMigrate(&models.TravelAdvice{})
	_ = DB.AutoMigrate(&models.TravelAdvice{})
	if err != nil {
		log.Fatalf("❌ Failed to auto-migrate: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
