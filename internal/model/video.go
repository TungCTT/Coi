package model

import "time"

type Video struct {
	ID                int               `json:"id" gorm:"primaryKey"`
	ChannelID         int               `json:"channel_id" gorm:"not null"`
	Channel           Channel           `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryID        int               `json:"category_id" gorm:"not null"`
	Category          Category          `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	VideoURL          string            `json:"video_url" gorm:"size:500"`
	ThumbnailURL      string            `json:"thumbnail_url" gorm:"size:500"`
	OriginalFileName  string            `json:"original_file_name"`
	FileSize          int64             `json:"file_size"`
	DurationSeconds   int64             `json:"duration"`
	Status            string            `json:"status" gorm:"type:varchar(30);not null;default:'processing'"`
	Visibility        VideoVisibility   `json:"visibility" gorm:"type:varchar(30);not null;default:'public'"`
	CommentsEnabled   bool              `json:"comments_enabled" gorm:"default:true"`
	ViewCount         int64             `json:"view_count" gorm:"default:0"`
	LikeCount         int64             `json:"like_count" gorm:"default:0"`
	DislikeCount      int64             `json:"dislike_count" gorm:"default:0"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	Comments          []Comment         `json:"comments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoReactions    []VideoReaction   `json:"video_reactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoTags         []VideoTag        `json:"video_tags" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlaylistItems     []PlaylistItem    `json:"playlist_items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ViewHistories     []ViewHistory     `json:"view_histories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}