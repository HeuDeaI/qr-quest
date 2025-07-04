package handlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"qr-quest/internal/repositories"
)

const hashedAdminPassword = "$2a$10$eD2aU6ZPEvcgF0//FIJl/uNvggYY5POOekaEsNxIDm61x2zxyHRzi"

type AdminHandler struct {
	questionRepository repositories.QuestionRepository
}

func NewAdminHandler(questionRepository *repositories.QuestionRepository) *AdminHandler {
	return &AdminHandler{
		questionRepository: *questionRepository,
	}
}

func (h *AdminHandler) ShowAdminLoginPage(c *gin.Context) {
	session := sessions.Default(c)

	errMsg := verifyLoginAvailability(session)

	c.HTML(http.StatusOK, "admin_login.html", errMsg)
}

func verifyLoginAvailability(session sessions.Session) gin.H {
	lastFailed, _ := session.Get("lastFailedAttempt").(int64)
	if time.Since(time.Unix(lastFailed, 0)) > 5*time.Minute {
		session.Delete("loginAttempts")
		session.Delete("lastFailedAttempt")
		session.Delete("lockedOut")
		session.Save()
		return gin.H{}
	}

	attempts, _ := session.Get("loginAttempts").(int)
	if attempts >= 3 {
		session.Set("lockedOut", true)
		session.Save()
		return gin.H{
			"error":     "Доступ временно заблокирован после 3 неверных попыток.",
			"lockedOut": true,
		}
	}

	flashes := session.Flashes()
	session.Save()

	var errMsg string
	if len(flashes) > 0 {
		errMsg = flashes[0].(string)
	} else {
		return gin.H{}
	}
	return gin.H{
		"error": errMsg,
	}
}

func (h *AdminHandler) HandleAdminLogin(c *gin.Context) {
	session := sessions.Default(c)
	if locked, ok := session.Get("lockedOut").(bool); ok && locked {
		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedAdminPassword), []byte(password))
	if err != nil {
		attempts, _ := session.Get("loginAttempts").(int)
		session.Set("lastFailedAttempt", time.Now().Unix())
		session.Set("loginAttempts", attempts+1)

		session.AddFlash("Неверный пароль. Попробуйте снова.")
		session.Save()

		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	session.Set("isAdmin", true)
	session.Delete("loginAttempts")
	session.Delete("lastFailedAttempt")
	session.Delete("lockedOut")
	session.Save()

	c.Redirect(http.StatusFound, "/admin/questions/list")
}

func (h *AdminHandler) ShowListOfQuestions(c *gin.Context) {
	listOfQuestions, err := h.questionRepository.GetListOfQuestions()
	if err != nil {

	}
	c.HTML(http.StatusOK, "list_of_questions.html", gin.H{
		"listOfQuestion": listOfQuestions,
	})
}

func (h *AdminHandler) ShowQuestionByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid UUID format")
		return
	}

	questionData, err := h.questionRepository.GetQuestionByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve question")
		return
	}

	c.HTML(http.StatusOK, "question_data.html", gin.H{
		"questionData": questionData,
	})
}

func (h *AdminHandler) HandleCreateQuestion(c *gin.Context) {
	// Implement your question creation logic here
	c.JSON(http.StatusOK, gin.H{"message": "Question created"})
}

func (h *AdminHandler) HandleListUsers(c *gin.Context) {
	// Implement your user listing logic here
	c.JSON(http.StatusOK, gin.H{"users": []string{"user1", "user2"}})
}
