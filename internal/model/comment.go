package model

import "time"

type Comment struct {
	ID               int                `json:"id" gorm:"primaryKey"`
	VideoID          int                `json:"video_id" gorm:"not null"`
	Video            Video              `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID           int                `json:"user_id" gorm:"not null"`
	User             User               `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentID         *int               `json:"parent_id"`
	Content          string             `json:"content"`
	Status           CommentStatus      `json:"status" gorm:"type:varchar(30);not null;default:'visible'"`
	Likes            int                `json:"likes" gorm:"default:0"`
	CreatedAt        time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
	IsUpdated        bool               `json:"is_updated" gorm:"default:false"`
	// Relationships
	CommentReactions []CommentReaction  `json:"comment_reactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}