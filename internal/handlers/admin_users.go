package handlers

import (
	"net/http"
	"strconv"

	"qr-quest/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ShowListOfUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении списка пользователей")
		return
	}
	c.HTML(http.StatusOK, "list_of_users.html", gin.H{
		"listOfUsers": users,
	})
}

func (h *AdminHandler) ShowUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить пользователя")
		return
	}
	c.HTML(http.StatusOK, "user_data.html", gin.H{
		"userData": user,
	})
}

func (h *AdminHandler) ShowEditUserPage(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить пользователя")
		return
	}
	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"userData": user,
	})
}

func (h *AdminHandler) HandleEditUser(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name   string `form:"name" binding:"required"`
		Points string `form:"points" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.String(http.StatusBadRequest, "Ошибка валидации: %v", err)
		return
	}

	pointsUint, err := strconv.ParseUint(input.Points, 10, 32)
	if err != nil {
		c.String(http.StatusBadRequest, "Очки должны быть числом")
		return
	}

	if err := h.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(models.User{Name: input.Name, Points: int(pointsUint)}).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при обновлении пользователя: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/admin/users/"+id)
}

func (h *AdminHandler) HandleDeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при удалении пользователя")
		return
	}
	c.Redirect(http.StatusFound, "/admin/users/list")
}
