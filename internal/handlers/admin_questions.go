package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"qr-quest/internal/models"
)

func (h *AdminHandler) ShowListOfQuestions(c *gin.Context) {
	var questions []models.Question
	if err := h.db.Find(&questions).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении списка вопросов")
		return
	}
	c.HTML(http.StatusOK, "list_of_questions.html", gin.H{
		"listOfQuestion": questions,
	})
}

func (h *AdminHandler) ShowQuestionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Некорректный UUID")
		return
	}

	var question models.Question
	if err := h.db.First(&question, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить вопрос")
		return
	}

	c.HTML(http.StatusOK, "question_data.html", gin.H{
		"questionData": question,
	})
}

func (h *AdminHandler) ShowCreateQuestionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_question.html", nil)
}

func (h *AdminHandler) HandleCreateQuestion(c *gin.Context) {
	var input struct {
		Text   string `form:"text" binding:"required"`
		Answer string `form:"answer" binding:"required"`
		Note   string `form:"note"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Ошибка валидации: %v", err)
		return
	}

	question := models.Question{
		ID:     uuid.New(),
		Text:   input.Text,
		Answer: input.Answer,
		Note:   input.Note,
	}

	if err := h.db.Create(&question).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при сохранении вопроса: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions/"+question.ID.String())
}

func (h *AdminHandler) ShowEditQuestionPage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Некорректный UUID")
		return
	}

	var question models.Question
	if err := h.db.First(&question, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить вопрос")
		return
	}

	c.HTML(http.StatusOK, "edit_question.html", gin.H{
		"questionData": question,
	})
}

func (h *AdminHandler) HandleEditQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Некорректный UUID")
		return
	}

	var input struct {
		Text   string `form:"text" binding:"required"`
		Answer string `form:"answer" binding:"required"`
		Note   string `form:"note"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Ошибка валидации: %v", err)
		return
	}

	if err := h.db.Model(&models.Question{}).
		Where("id = ?", id).
		Updates(models.Question{Text: input.Text, Answer: input.Answer, Note: input.Note}).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при обновлении вопроса: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions/"+id.String())
}

func (h *AdminHandler) HandleDeleteQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Некорректный UUID")
		return
	}

	if err := h.db.Delete(&models.Question{}, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при удалении вопроса")
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions/list")
}
