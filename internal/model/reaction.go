package model

import "time"

type VideoReaction struct {
	ID        int          `json:"id" gorm:"primaryKey"`
	UserID    int          `json:"user_id" gorm:"not null;index:idx_video_reaction_user_video,unique"`
	User      User         `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID   int          `json:"video_id" gorm:"not null;index:idx_video_reaction_user_video,unique"`
	Video     Video        `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      ReactionType `json:"type" gorm:"type:varchar(30);not null"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

func (VideoReaction) TableName() string {
	return "video_reactions"
}