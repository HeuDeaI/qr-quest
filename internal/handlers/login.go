package handlers

import (
	"net/http"
	"qr-quest/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *UserHandler) HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Имя пользователя обязательно",
		})
		return
	}

	var existing models.User
	if err := h.db.Where("name = ?", username).First(&existing).Error; err == nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Имя уже занято. Пожалуйста, выберите другое",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.String(http.StatusInternalServerError, "Ошибка проверки имени пользователя")
		return
	}

	newUser := models.User{
		Name:   username,
		Points: 0,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при создании пользователя")
		return
	}

	session := sessions.Default(c)
	session.Set("Username", newUser.Name)
	session.Save()

	c.Redirect(http.StatusFound, "/about")
}

func (h *UserHandler) ShowAboutPage(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("Username")
	c.HTML(http.StatusOK, "about.html", gin.H{"Username": username})
}
