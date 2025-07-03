package models

import (
	"github.com/google/uuid"
)

type Question struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Text   string
	Answer string
}
