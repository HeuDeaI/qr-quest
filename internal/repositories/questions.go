package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"qr-quest/internal/models"
)

type QuestionRepository interface {
	GetListOfQuestions() ([]models.Question, error)
	GetQuestionByID(id uuid.UUID) (*models.Question, error)
	CreateQuestion(question *models.Question) error
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(id uuid.UUID) error
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

func (repo *questionRepository) GetQuestionByID(id uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := repo.db.First(&question, "id = ?", id).Error
	return &question, err
}

func (repo *questionRepository) CreateQuestion(question *models.Question) error {
	return repo.db.Create(question).Error
}

func (repo *questionRepository) UpdateQuestion(question *models.Question) error {
	return repo.db.Save(question).Error
}

func (repo *questionRepository) DeleteQuestion(id uuid.UUID) error {
	return repo.db.Delete(&models.Question{}, "id = ?", id).Error
}
