package main

import (
	"qr-quest/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres dbname=qr_quest sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	router := gin.Default()

	server.SetupRouter(router, db)
	router.Run(":8080")
}
