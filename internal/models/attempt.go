package models

import (
	"github.com/google/uuid"
)

type UserQuestionAttempt struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	UserID     uint
	QuestionID uuid.UUID
	Attempts   uint
	Correct    bool
}
