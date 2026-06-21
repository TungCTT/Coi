package main

import (
	"log"

	"coi/config"
	categoryRoutes "coi/internal/category"
	channelRoutes "coi/internal/channel"
	"coi/internal/model"
	userRoutes "coi/internal/user"
	videoRoutes "coi/internal/video"
	"coi/pkg/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbCfg := config.LoadDBConfig()
	s3Cfg, err := config.LoadS3Config()
	if err != nil {
		log.Fatalf("Khong the load R2 config: %v", err)
	}

	db, err := gorm.Open(mysql.Open(dbCfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Khong the ket noi DB: %v", err)
	}
	log.Println("Ket noi database thanh cong")

	r2Client, err := storage.NewR2Client(s3Cfg.R2Endpoint, s3Cfg.R2AccessKeyID, s3Cfg.R2SecretAccessKey)
	if err != nil {
		log.Fatalf("Khong the ket noi R2: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Channel{}, &model.Category{}, &model.Video{}); err != nil {
		log.Fatalf("AutoMigrate that bai: %v", err)
	}
	log.Println("AutoMigrate thanh cong")

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		userRoutes.RegisterRoutes(v1, db)
		channelRoutes.RegisterRoutes(v1, db)
		categoryRoutes.RegisterRoutes(v1, db)
		videoRoutes.RegisterRoutes(v1, db, r2Client, s3Cfg.R2BucketName, s3Cfg.R2PublicBaseURL)
	}

	log.Println("Server dang chay tai http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Khong the khoi dong server: %v", err)
	}
}
