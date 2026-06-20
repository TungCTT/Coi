package model

import "time"

type Channel struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id" gorm:"uniqueIndex;not null"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string    `json:"name"`
	Handle      string    `json:"handle" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url" gorm:"size:500"`
	BannerURL   string    `json:"banner_url" gorm:"size:500"`
	Subscribers int       `json:"subscribers" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	Videos        []Video        `json:"videos" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateChannelRequest struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name" binding:"required"`
	Handle      string `json:"handle" binding:"required"`
	Description string `json:"description" binding:"required"`
	AvatarURL   string `json:"avatar_url" binding:"required"`
	BannerURL   string `json:"banner_url" binding:"required"`
}
