package user

import (
	"coi/internal/middleware"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, r2Client *s3.Client, bucketName string, publicBaseURL string) {
	repo := NewUserRepository(db)
	mediaStorage := NewR2UserMediaStorage(r2Client, bucketName, publicBaseURL)
	service := NewAuthService(repo, mediaStorage)
	handler := NewAuthHandler(service)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)

		protected := auth.Group("", middleware.AuthMiddleware())
		{
			protected.GET("/me", handler.Me)
			protected.POST("/me/media-upload-sessions", handler.CreateMediaUploadSession)
			protected.POST("/me/media/confirm-upload", handler.ConfirmMediaUpload)
		}
	}
}
