package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"qr-quest/internal/middlewares"
)

const hashedAdminPassword = "$2a$10$eD2aU6ZPEvcgF0//FIJl/uNvggYY5POOekaEsNxIDm61x2zxyHRzi"

func RegisterAdminRoutes(router *gin.Engine) {
	store := cookie.NewStore([]byte("super-secret-key"))
	router.Use(sessions.Sessions("mysession", store))

	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/login", handleAdminLogin)
	}

	protectedGroup := adminGroup.Group("/", middlewares.RequireAdminSession())

	questionsGroup := protectedGroup.Group("/questions")
	{
		questionsGroup.POST("/create", handleCreateQuestion)
	}

	usersGroup := protectedGroup.Group("/users")
	{
		usersGroup.GET("/all", handleListUsers)
	}
}

func handleAdminLogin(c *gin.Context) {
	password := c.PostForm("password")

	err := bcrypt.CompareHashAndPassword([]byte(hashedAdminPassword), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	session := sessions.Default(c)
	session.Set("isAdmin", true)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Admin login successful"})
}

func handleCreateQuestion(c *gin.Context) {
	// Implement your question creation logic here
	c.JSON(http.StatusOK, gin.H{"message": "Question created"})
}

func handleListUsers(c *gin.Context) {
	// Implement your user listing logic here
	c.JSON(http.StatusOK, gin.H{"users": []string{"user1", "user2"}})
}
