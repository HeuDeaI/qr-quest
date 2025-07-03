package main

import (
	"github.com/gin-gonic/gin"
	"qr-quest/internal/handlers"
)

func main() {
	router := gin.Default()

	handlers.RegisterAdminRoutes(router)
}
