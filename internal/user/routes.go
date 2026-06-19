package user

import (
	"coi/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes "wire" tất cả các layers lại với nhau và đăng ký routes vào Gin router.
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Khởi tạo từng layer theo thứ tự dependency (từ dưới lên trên)
	repo    := NewUserRepository(db)   
	service := NewAuthService(repo)   
	handler := NewAuthHandler(service) 

	auth := rg.Group("/auth")
	{
		// Public routes 
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)

		// Protected routes
		protected := auth.Group("", middleware.AuthMiddleware())
		{
			protected.GET("/me", handler.Me)
		}
	}
}
