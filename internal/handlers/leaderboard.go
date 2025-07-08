package handlers

import (
	"net/http"
	"qr-quest/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ShowLeaderBoard(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")

	var allUsers []models.User
	if err := h.db.Order("points DESC").Find(&allUsers).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки таблицы лидеров")
		return
	}

	c.HTML(http.StatusOK, "leaderboard.html", gin.H{
		"User":     user,
		"AllUsers": allUsers,
	})
}
