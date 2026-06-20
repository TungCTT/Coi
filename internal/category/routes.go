package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewCategoryRepository(db)
	service := NewCategoryService(repo)
	handler := NewCategoryHandler(service)

	categories := rg.Group("/categories")
	{
		categories.GET("", handler.ListCategories)
		categories.GET("/:id", handler.GetCategoryByID)
	}
}
