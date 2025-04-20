package models

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	TravelAdviceID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt     time.Time
}
