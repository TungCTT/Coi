package model

import "time"

type PlaylistItem struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	PlaylistID int       `json:"playlist_id" gorm:"not null;index:idx_playlist_item_playlist_video,unique"`
	Playlist   Playlist  `json:"playlist" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID    int       `json:"video_id" gorm:"not null;index:idx_playlist_item_playlist_video,unique"`
	Video      Video     `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SortOrder  int       `json:"sort_order" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}