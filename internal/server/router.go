package server

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qr-quest/internal/handlers"
	"qr-quest/internal/middlewares"
	"qr-quest/internal/models"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	funcMap := template.FuncMap{
		"add": func(i, j int) int {
			return i + j
		},
	}

	router.SetFuncMap(funcMap)
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	db.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.UserQuestionAttempt{},
	)

	adminHandler := handlers.NewAdminHandler(db)
	userHandler := handlers.NewUserHandler(db)

	RegisterAdminRoutes(router, adminHandler, userHandler)
}

func RegisterAdminRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler, userHandler *handlers.UserHandler) {
	store := cookie.NewStore([]byte("your-secret-key"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   10 * 24 * 60 * 60,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("static/favicon.ico")
	})

	router.GET("/login", userHandler.ShowLoginPage)
	router.POST("/login", userHandler.HandleLogin)
	router.GET("/about", userHandler.ShowAboutPage)
	router.GET("/leaderboard", userHandler.ShowLeaderBoard)

	questGroup := router.Group("/questions", middlewares.RequireUserSession())
	{
		questGroup.GET("/:id", userHandler.ShowQuestion)
		questGroup.POST("/:id", userHandler.SubmitAnswer)
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", adminHandler.ShowAdminLoginPage)
		adminGroup.POST("/login", adminHandler.HandleAdminLogin)
	}

	protectedGroup := adminGroup.Group("/", middlewares.RequireAdminSession())

	protectedGroup.GET("/home", adminHandler.ShowAdminHomePage)

	questionsGroup := protectedGroup.Group("/questions")
	{
		questionsGroup.GET("/list", adminHandler.ShowListOfQuestions)
		questionsGroup.GET("/:id", adminHandler.ShowQuestionByID)
		questionsGroup.GET("/create", adminHandler.ShowCreateQuestionPage)
		questionsGroup.POST("/create", adminHandler.HandleCreateQuestion)
		questionsGroup.GET("/:id/qr", adminHandler.GenerateQRCodePDF)
		questionsGroup.POST("/:id/delete", adminHandler.HandleDeleteQuestion)
		questionsGroup.GET("/:id/edit", adminHandler.ShowEditQuestionPage)
		questionsGroup.POST("/:id/edit", adminHandler.HandleEditQuestion)
	}

	usersGroup := protectedGroup.Group("/users")
	{
		usersGroup.GET("/list", adminHandler.ShowListOfUsers)
		usersGroup.GET("/:id", adminHandler.ShowUserByID)
		usersGroup.GET("/:id/edit", adminHandler.ShowEditUserPage)
		usersGroup.POST("/:id/edit", adminHandler.HandleEditUser)
		usersGroup.POST("/:id/delete", adminHandler.HandleDeleteUser)
	}
}
