package channel

import (
	"coi/internal/middleware"
	"coi/internal/user"
	"coi/pkg/txmanager"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, r2Client *s3.Client, bucketName string, publicBaseURL string) {
	channelRepo := NewChannelRepository(db)
	userRepo := user.NewUserRepository(db)
	txManager := txmanager.NewTxManager(db)
	mediaStorage := NewR2ChannelMediaStorage(r2Client, bucketName, publicBaseURL)
	service := NewChannelService(channelRepo, userRepo, txManager, mediaStorage)
	handler := NewChannelHandler(service)

	channels := rg.Group("/channels")
	{
		channels.GET("/:id", handler.GetChannelByID)
		channels.GET("", handler.GetChannelsByName)

		protected := channels.Group("", middleware.AuthMiddleware())
		{
			protected.POST("", handler.CreateChannel)
			protected.POST("/:id/media-upload-sessions", handler.CreateMediaUploadSession)
			protected.POST("/:id/media/confirm-upload", handler.ConfirmMediaUpload)
		}
	}

	me := rg.Group("/me", middleware.AuthMiddleware())
	{
		me.GET("/channel", handler.GetMyChannel)
	}
}
