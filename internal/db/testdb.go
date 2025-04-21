package db

import (
	"log"
	"github.com/thoraf20/smart-travel-advisor/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func InitTestDB() {
	var err error

	TestDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open test SQLite DB: %v", err)
	}

	err = TestDB.AutoMigrate(
		&models.User{},
		&models.TravelAdvice{},
		&models.City{},
		&models.Favorite{},
		&models.UserPreferences{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate test DB: %v", err)
	}
}
