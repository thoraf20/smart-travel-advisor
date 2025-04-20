package models

import (
	"time"

	"github.com/google/uuid"
)

type TravelAdvice struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	StartDate time.Time
	EndDate   time.Time
	Cities    []byte `gorm:"type:jsonb"`
	Advice    []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
}
