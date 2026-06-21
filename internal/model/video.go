package model

import "time"

type Video struct {
	ID               int             `json:"id" gorm:"primaryKey"`
	ChannelID        int             `json:"channel_id" gorm:"not null"`
	Channel          Channel         `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryID       int             `json:"category_id" gorm:"not null"`
	Category         Category        `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	StorageKey       string          `json:"-" gorm:"size:500;not null;uniqueIndex"`
	VideoURL         string          `json:"video_url" gorm:"size:500"`
	ThumbnailURL     string          `json:"thumbnail_url" gorm:"size:500"`
	OriginalFileName string          `json:"original_file_name"`
	FileSize         int64           `json:"file_size"`
	ContentType      string          `json:"content_type" gorm:"size:100"`
	ETag             string          `json:"etag" gorm:"size:200"`
	DurationSeconds  int64           `json:"duration"`
	Status           VideoStatus     `json:"status" gorm:"type:varchar(30);not null;default:'uploading'"`
	Visibility       VideoVisibility `json:"visibility" gorm:"type:varchar(30);not null;default:'public'"`
	CommentsEnabled  bool            `json:"comments_enabled" gorm:"default:true"`
	ViewCount        int64           `json:"view_count" gorm:"default:0"`
	LikeCount        int64           `json:"like_count" gorm:"default:0"`
	DislikeCount     int64           `json:"dislike_count" gorm:"default:0"`
	UploadExpiresAt  time.Time       `json:"upload_expires_at"`
	UploadedAt       *time.Time      `json:"uploaded_at"`
	CreatedAt        time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	Comments       []Comment       `json:"comments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoReactions []VideoReaction `json:"video_reactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoTags      []VideoTag      `json:"video_tags" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlaylistItems  []PlaylistItem  `json:"playlist_items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ViewHistories  []ViewHistory   `json:"view_histories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateVideoRequest struct {
	ChannelID   int             `json:"channel_id" binding:"required"`
	CategoryID  int             `json:"category_id" binding:"required"`
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description"`
	Visibility  VideoVisibility `json:"visibility" binding:"required,oneof=public private"`
}

type CreateVideoUploadRequest struct {
	ChannelID        int             `json:"channel_id" binding:"required"`
	CategoryID       int             `json:"category_id" binding:"required"`
	Title            string          `json:"title" binding:"required"`
	Description      string          `json:"description"`
	Visibility       VideoVisibility `json:"visibility" binding:"required,oneof=public private"`
	OriginalFileName string          `json:"original_file_name" binding:"required"`
	FileSize         int64           `json:"file_size" binding:"required"`
	ContentType      string          `json:"content_type" binding:"required"`
}

type CreateVideoUploadResponse struct {
	VideoID         int       `json:"video_id"`
	UploadURL       string    `json:"upload_url"`
	StorageKey      string    `json:"storage_key"`
	Status          string    `json:"status"`
	UploadExpiresAt time.Time `json:"upload_expires_at"`
}

type ConfirmVideoUploadRequest struct {
	VideoID int `json:"video_id" binding:"required"`
}
