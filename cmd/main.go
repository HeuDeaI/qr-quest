package main

import (
	"context"
	"qr-quest/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db port=5432 user=postgres password=postgres dbname=qr_quest sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("Failed to get generic DB interface: %v", err)
	}
	_, err = sqlDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		logrus.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, "qrquest-images")
	if err != nil {
		logrus.Fatalf("Failed to check bucket: %v", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, "qrquest-images", minio.MakeBucketOptions{}); err != nil {
			logrus.Fatalf("Failed to create bucket: %v", err)
		}
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server.SetupRouter(router, db, minioClient)
	router.Run(":8080")
}
