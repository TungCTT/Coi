package video

import (
	"coi/internal/category"
	"coi/internal/channel"
	"coi/internal/middleware"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, r2Client *s3.Client, bucketName string, publicBaseURL string) {
	videoRepo := NewVideoRepository(db)
	videoStorage := NewR2VideoStorage(r2Client, bucketName, publicBaseURL)
	channelRepo := channel.NewChannelRepository(db)
	categoryRepo := category.NewCategoryRepository(db)
	service := NewVideoService(videoRepo, videoStorage, channelRepo, categoryRepo)
	handler := NewVideoHandler(service)

	videos := rg.Group("/videos")
	{
		videos.GET("", handler.ListVideos)
		videos.GET("/:id", handler.GetVideoByID)
		videos.GET("/:id/stream", handler.StreamVideo)

		protected := videos.Group("", middleware.AuthMiddleware())
		{
			protected.POST("/upload-sessions", handler.CreateUploadSession)
			protected.POST("/:id/confirm-upload", handler.ConfirmUpload)
			protected.DELETE("/:id", handler.DeleteVideo)
		}
	}

	rg.GET("/channels/:id/videos", handler.ListVideosByChannel)

	internal := rg.Group("/internal/videos")
	{
		internal.POST("/upload-completed", handler.HandleUploadCompleted)
	}
}
