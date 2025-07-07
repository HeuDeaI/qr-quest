package handlers

import (
	"fmt"
	"math"
	"net/http"
	"qr-quest/internal/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ShowQuestion(c *gin.Context) {
	questionID := c.Param("id")

	var question models.Question
	if err := h.db.First(&question, "id = ?", questionID).Error; err != nil {
		c.String(http.StatusNotFound, "Вопрос не найден")
		return
	}

	user, err := h.getCurrentUser(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	var attempt models.UserQuestionAttempt
	err = h.db.
		Where("user_id = ? AND question_id = ?", user.ID, question.ID).
		First(&attempt).Error

	if err == nil && attempt.Correct {
		c.HTML(http.StatusForbidden, "already_answered.html", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "question.html", gin.H{
		"question": question,
	})
}

func (h *UserHandler) SubmitAnswer(c *gin.Context) {
	questionID := c.Param("id")
	userAnswer := c.PostForm("answer")

	var question models.Question
	if err := h.db.First(&question, "id = ?", questionID).Error; err != nil {
		c.String(http.StatusNotFound, "Вопрос не найден")
		return
	}

	user, err := h.getCurrentUser(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	isCorrect := question.Answer == userAnswer

	var attempt models.UserQuestionAttempt
	err = h.db.
		Where("user_id = ? AND question_id = ?", user.ID, question.ID).
		First(&attempt).Error

	if attempt.Correct {
		h.ShowQuestion(c)
		return
	}

	earnedPoints := 0
	if err != nil && err.Error() == "record not found" {
		attempt = models.UserQuestionAttempt{
			UserID:     user.ID,
			QuestionID: question.ID,
			Attempts:   1,
			Correct:    isCorrect,
		}
		if isCorrect {
			earnedPoints = h.calculatePointsForCorrectAnswer(question, attempt)
			user.Points += earnedPoints
		}
		h.db.Create(&attempt)
	} else {
		attempt.Attempts++
		if isCorrect && !attempt.Correct {
			attempt.Correct = true
			earnedPoints = h.calculatePointsForCorrectAnswer(question, attempt)
			user.Points += earnedPoints
		}
		h.db.Save(&attempt)
	}

	h.db.Save(&user)

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Correct":      isCorrect,
		"Question":     question,
		"Answer":       userAnswer,
		"Points":       user.Points,
		"EarnedPoints": earnedPoints,
	})
}

func (h *UserHandler) calculatePointsForCorrectAnswer(question models.Question, attempt models.UserQuestionAttempt) int {
	elapsed := float64(time.Now().Unix()-question.CreatedAt) / 100
	if elapsed < 1 {
		elapsed = 1
	}

	timePenalty := float64(question.Points)*0.2 - float64(question.Points)*0.05*math.Log(elapsed)

	score := float64(question.Points)*math.Exp(-float64(attempt.Attempts-1)/10) + timePenalty

	if score < float64(question.Points)*0.5 {
		score = float64(question.Points) * 0.5
	}

	return int(math.Round(score))
}

func (h *UserHandler) getCurrentUser(c *gin.Context) (*models.User, error) {
	session := sessions.Default(c)
	username := session.Get("username")

	if username == nil {
		return nil, fmt.Errorf("пользователь не найден в сессии")
	}

	var user models.User
	if err := h.db.Where("name = ?", username.(string)).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
