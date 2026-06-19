package model

import "time"

type Playlist struct {
	ID            int              `json:"id" gorm:"primaryKey"`
	UserID        int              `json:"user_id" gorm:"not null;index"`
	User          User             `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name          string           `json:"name" gorm:"not null"`
	Description   string           `json:"description"`
	IsPublic      bool             `json:"is_public" gorm:"default:false"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	PlaylistItems []PlaylistItem   `json:"playlist_items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}