package model

import "time"

type Tag struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:120;uniqueIndex;not null"`
	Slug      string    `json:"slug" gorm:"size:160;uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	VideoTags []VideoTag `json:"video_tags" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
