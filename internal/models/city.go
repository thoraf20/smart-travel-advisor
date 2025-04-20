package models

import (
	"github.com/google/uuid"
)

type City struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Country     string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Latitude    float64
	Longitude   float64
}
