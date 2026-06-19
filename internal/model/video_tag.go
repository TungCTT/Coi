package model

import "time"

type VideoTag struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	VideoID   int       `json:"video_id" gorm:"not null;index:idx_video_tag_video_tag,unique"`
	TagID     int       `json:"tag_id" gorm:"not null;index:idx_video_tag_video_tag,unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Video     Video     `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag       `json:"tag" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}