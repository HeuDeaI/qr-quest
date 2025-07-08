package main

import (
	"qr-quest/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=postgres-db port=5432 user=postgres password=postgres dbname=qr_quest sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("Failed to get generic DB interface: %v", err)
	}
	_, err = sqlDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		logrus.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server.SetupRouter(router, db)
	router.Run(":8080")
}
