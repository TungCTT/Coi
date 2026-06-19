package model

import "time"

type CommentReaction struct {
	ID        int          `json:"id" gorm:"primaryKey"`
	UserID    int          `json:"user_id" gorm:"not null;index:idx_comment_reaction_user_comment,unique"`
	User      User         `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentID int          `json:"comment_id" gorm:"not null;index:idx_comment_reaction_user_comment,unique"`
	Comment   Comment      `json:"comment" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      ReactionType `json:"type" gorm:"type:varchar(30);not null"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

func (CommentReaction) TableName() string {
	return "comment_reactions"
}