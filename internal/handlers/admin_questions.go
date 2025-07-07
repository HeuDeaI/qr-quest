package handlers

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/phpdave11/gofpdf"
	"github.com/skip2/go-qrcode"

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
		Points int    `form:"points" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Ошибка валидации: %v", err)
		return
	}
	question := models.Question{
		ID:        uuid.New(),
		Text:      input.Text,
		Answer:    input.Answer,
		Note:      input.Note,
		Points:    input.Points,
		CreatedAt: time.Now().Unix(),
	}

	if err := h.db.Create(&question).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при сохранении вопроса: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions/"+question.ID.String())
}

func (h *AdminHandler) GenerateQRCodePDF(c *gin.Context) {
	id := c.Param("id")

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + c.Request.Host + "/questions/" + id

	qrPNG, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка генерации QR-кода")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("ComicSansMS", "", "internal/handlers/ComicSansMS.ttf")
	pdf.SetFont("ComicSansMS", "", 72)
	pdf.AddPage()
	text := "QR-Квест!"
	width, _ := pdf.GetPageSize()
	textWidth := pdf.GetStringWidth(text)
	pdf.SetXY((width-textWidth)/2, 20)
	pdf.CellFormat(textWidth, 10, text, "", 0, "C", false, 0, "")

	imgOpts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
	pdf.RegisterImageOptionsReader("qr.png", imgOpts, bytes.NewReader(qrPNG))

	pdf.ImageOptions("qr.png", 15, 40, 180, 180, false, imgOpts, 0, "")

	pdf.SetFont("ComicSansMS", "", 48)
	pdf.SetXY(15, 230)
	pdf.MultiCell(180, 20, "Сканируй и отвечай на вопросы!", "", "C", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка генерации PDF-файла")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=qr-question-"+id+".pdf")
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
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
		Points int    `form:"points" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Ошибка валидации: %v", err)
		return
	}

	if err := h.db.Model(&models.Question{}).
		Where("id = ?", id).
		Updates(models.Question{Text: input.Text, Answer: input.Answer, Note: input.Note, Points: input.Points}).Error; err != nil {
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
