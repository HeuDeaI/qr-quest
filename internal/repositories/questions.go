package repositories

import (
	"gorm.io/gorm"
	"qr-quest/internal/models"
)

type QuestionRepository interface {
	GetListOfQuestions() ([]models.Question, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (repo *questionRepository) GetListOfQuestions() ([]models.Question, error) {
	var questions []models.Question
	err := repo.db.Find(&questions).Error
	return questions, err
}
