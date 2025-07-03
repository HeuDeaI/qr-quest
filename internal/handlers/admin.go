package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	c.HTML(http.StatusOK, "admin_login.html", gin.H{})
}

func (h *AdminHandler) HandleAdminLogin(c *gin.Context) {
	password := c.PostForm("password")

	err := bcrypt.CompareHashAndPassword([]byte(hashedAdminPassword), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	session := sessions.Default(c)
	session.Set("isAdmin", true)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Admin login successful"})
}

func (h *AdminHandler) ShowListOfQuestions(c *gin.Context) {
	listOfQuestions, err := h.questionRepository.GetListOfQuestions()
	if err != nil {

	}
	c.HTML(http.StatusOK, "list_of_questions.html", listOfQuestions)
}

func (h *AdminHandler) HandleCreateQuestion(c *gin.Context) {
	// Implement your question creation logic here
	c.JSON(http.StatusOK, gin.H{"message": "Question created"})
}

func (h *AdminHandler) HandleListUsers(c *gin.Context) {
	// Implement your user listing logic here
	c.JSON(http.StatusOK, gin.H{"users": []string{"user1", "user2"}})
}
