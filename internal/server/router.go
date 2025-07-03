package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qr-quest/internal/handlers"
	"qr-quest/internal/middlewares"
	"qr-quest/internal/repositories"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	router.LoadHTMLGlob("web/templates/*")

	questionRepository := repositories.NewQuestionRepository(db)
	adminHandler := handlers.NewAdminHandler(&questionRepository)

	RegisterAdminRoutes(router, adminHandler)
}

func RegisterAdminRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler) {
	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("mysession", store))

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", adminHandler.ShowAdminLoginPage)
		adminGroup.POST("/login", adminHandler.HandleAdminLogin)
	}

	protectedGroup := adminGroup.Group("/", middlewares.RequireAdminSession())

	questionsGroup := protectedGroup.Group("/questions")
	{
		questionsGroup.GET("/list", adminHandler.ShowListOfQuestions)
		// questionsGroup.GET("/:uid", adminHandler.ShowQuestion)
		questionsGroup.POST("/create", adminHandler.HandleCreateQuestion)
	}

	usersGroup := protectedGroup.Group("/users")
	{
		usersGroup.GET("/list", adminHandler.HandleListUsers)
	}
}
