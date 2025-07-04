package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if _, ok := session.Get("username").(string); !ok {
			session.Set("redirectTo", c.Request.URL.Path)
			session.Save()

			c.Redirect(http.StatusFound, "/about")
			c.Abort()
			return
		}
		c.Next()
	}
}
