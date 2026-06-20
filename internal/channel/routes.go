package channel

import (
	"coi/internal/middleware"
	"coi/internal/user"
	"coi/pkg/txmanager"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	channelRepo := NewChannelRepository(db)
	userRepo := user.NewUserRepository(db)
	txManager := txmanager.NewTxManager(db)
	service := NewChannelService(channelRepo, userRepo, txManager)
	handler := NewChannelHandler(service)

	channels := rg.Group("/channels")
	{
		channels.GET("/:id", handler.GetChannelByID)
		channels.GET("", handler.GetChannelsByName)

		protected := channels.Group("", middleware.AuthMiddleware())
		{
			protected.POST("", handler.CreateChannel)
		}
	}
}
