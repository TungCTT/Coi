package main

import (
	"log"

	"coi/config"
	categoryRoutes "coi/internal/category"
	channelRoutes "coi/internal/channel"
	"coi/internal/model"
	userRoutes "coi/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbCfg := config.LoadDBConfig()

	//Kết nối database qua GORM
	db, err := gorm.Open(mysql.Open(dbCfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}
	log.Println("Kết nối database thành công")

	//Auto migrate
	//hình như nên dùng go-migrate
	if err := db.AutoMigrate(&model.User{}, &model.Channel{}, &model.Category{}); err != nil {
		log.Fatalf("AutoMigrate thất bại: %v", err)
	}
	log.Println("AutoMigrate thành công")

	//Khởi tạo Gin engine
	//tạo engine với 2 middleware mặc định: Logger và Recovery.

	r := gin.Default()

	//Đăng ký routes
	v1 := r.Group("/api/v1")
	{
		userRoutes.RegisterRoutes(v1, db)
		channelRoutes.RegisterRoutes(v1, db)
		categoryRoutes.RegisterRoutes(v1, db)
	}

	log.Println("Server đang chạy tại http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
