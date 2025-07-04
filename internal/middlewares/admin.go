package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAdminSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if isAdmin, ok := session.Get("isAdmin").(bool); !ok || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
			c.Abort()
			return
		}
		c.Next()
	}
}
