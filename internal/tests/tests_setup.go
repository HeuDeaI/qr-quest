package tests

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"qr-quest/internal/handlers"
)

func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	store := cookie.NewStore([]byte("test-secret"))
	router.Use(sessions.Sessions("mysession", store))

	handlers.RegisterAdminRoutes(router)
	return router
}
