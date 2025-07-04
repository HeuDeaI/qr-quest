package handlers

import (
	"net/http"
	"qr-quest/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ShowQuestion(c *gin.Context) {
	questionID := c.Param("id")

	var question models.Question
	if err := h.db.First(&question, "id = ?", questionID).Error; err != nil {
		c.String(http.StatusNotFound, "Вопрос не найден")
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

	isCorrect := question.Answer == userAnswer

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Correct":  isCorrect,
		"Question": question,
		"Answer":   userAnswer,
	})
}
