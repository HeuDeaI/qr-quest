package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name             string
	Points           uint
	TriedQuestionIDs []uint     `gorm:"-"`
	TriedQuestions   []Question `gorm:"many2many:user_tried_questions;"`
}
