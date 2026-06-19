package model

import "time"

type Category struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url" gorm:"size:500"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}