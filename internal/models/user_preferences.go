package models

import (
	"github.com/google/uuid"
)

type UserPreferences struct {
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Preferences []byte   `gorm:"type:jsonb"` // stored as raw JSON
}
