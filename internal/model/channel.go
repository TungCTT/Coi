package model

import "time"

type Channel struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UserID           int       `json:"user_id" gorm:"uniqueIndex;not null"`
	User             User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name             string    `json:"name" gorm:"size:120"`
	Handle           string    `json:"handle" gorm:"size:120;uniqueIndex;not null"`
	Description      string    `json:"description"`
	AvatarURL        string    `json:"avatar_url" gorm:"size:500"`
	AvatarStorageKey string    `json:"-" gorm:"size:500"`
	BannerURL        string    `json:"banner_url" gorm:"size:500"`
	BannerStorageKey string    `json:"-" gorm:"size:500"`
	Subscribers      int       `json:"subscribers" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	Videos        []Video        `json:"videos" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateChannelRequest struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name" binding:"required"`
	Handle      string `json:"handle" binding:"required"`
	Description string `json:"description" binding:"required"`
	AvatarURL   string `json:"avatar_url"`
	BannerURL   string `json:"banner_url"`
}

type CreateChannelMediaUploadRequest struct {
	MediaType        string `json:"media_type" binding:"required,oneof=avatar banner"`
	OriginalFileName string `json:"original_file_name" binding:"required"`
	FileSize         int64  `json:"file_size" binding:"required"`
	ContentType      string `json:"content_type" binding:"required"`
}

type ConfirmChannelMediaUploadRequest struct {
	MediaType  string `json:"media_type" binding:"required,oneof=avatar banner"`
	StorageKey string `json:"storage_key" binding:"required"`
}

type ChannelMediaUploadResponse struct {
	UploadURL  string `json:"upload_url"`
	StorageKey string `json:"storage_key"`
	PublicURL  string `json:"public_url"`
}
