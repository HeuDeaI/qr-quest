package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"qr-quest/internal/handlers"
	"qr-quest/internal/middlewares"
	"qr-quest/internal/models"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	router.LoadHTMLGlob("templates/*")
	db.AutoMigrate(
		&models.User{},
		&models.Question{},
	)

	adminHandler := handlers.NewAdminHandler(db)

	RegisterAdminRoutes(router, adminHandler)
}

func RegisterAdminRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler) {
	store := cookie.NewStore([]byte("your-secret-key"))
	store.Options(sessions.Options{
		SameSite: http.SameSiteLaxMode, // Important for Safari
	})
	router.Use(sessions.Sessions("mysession", store))

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", adminHandler.ShowAdminLoginPage)
		adminGroup.POST("/login", adminHandler.HandleAdminLogin)
		adminGroup.GET("/home", adminHandler.ShowAdminHomePage)

	}

	protectedGroup := adminGroup.Group("/", middlewares.RequireAdminSession())

	questionsGroup := protectedGroup.Group("/questions")
	{
		questionsGroup.GET("/list", adminHandler.ShowListOfQuestions)
		questionsGroup.GET("/:id", adminHandler.ShowQuestionByID)
		questionsGroup.GET("/create", adminHandler.ShowCreateQuestionPage)
		questionsGroup.POST("/create", adminHandler.HandleCreateQuestion)
		questionsGroup.POST("/:id/delete", adminHandler.HandleDeleteQuestion)
		questionsGroup.GET("/:id/edit", adminHandler.ShowEditQuestionPage)
		questionsGroup.POST("/:id/edit", adminHandler.HandleEditQuestion)
	}

	// usersGroup := protectedGroup.Group("/users")
	// {
	// usersGroup.GET("/list", adminHandler.HandleListUsers)
	// }
}
