package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if _, ok := session.Get("Username").(string); !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
