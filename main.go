package main

import (
	"log"

	"coi/config"
	"coi/internal/model"
	userRoutes "coi/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ─── 1. Load cấu hình từ file .env ───────────────────────────────────────
	dbCfg := config.LoadDBConfig()

	// ─── 2. Kết nối database qua GORM ────────────────────────────────────────
	// GORM là ORM (Object-Relational Mapping) — map struct Go ↔ bảng DB.
	db, err := gorm.Open(mysql.Open(dbCfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}
	log.Println("Kết nối database thành công")

	// ─── 3. Auto migrate ─────────────────────────────────────────────────────
	// AutoMigrate tự động tạo / cập nhật bảng trong DB dựa theo struct.
	// Chỉ thêm cột mới và index mới — KHÔNG xóa cột cũ (an toàn).
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("AutoMigrate thất bại: %v", err)
	}
	log.Println("AutoMigrate thành công")

	// ─── 4. Khởi tạo Gin engine ──────────────────────────────────────────────
	// gin.Default() tạo engine với 2 middleware mặc định:
	//   - Logger: log mỗi request (method, path, status, latency)
	//   - Recovery: bắt panic và trả 500 thay vì crash server
	r := gin.Default()

	// ─── 5. Đăng ký routes ───────────────────────────────────────────────────
	// Tạo route group /api/v1 — tất cả routes sẽ có prefix này.
	// Versioning API (v1, v2...) là best practice để không breaking change client cũ.
	v1 := r.Group("/api/v1")
	{
		userRoutes.RegisterRoutes(v1, db)
	}

	// ─── 6. Khởi động server ─────────────────────────────────────────────────
	log.Println("Server đang chạy tại http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}