package model

import "time"

type User struct {
	ID              int              `json:"id" gorm:"primaryKey"`
	Name            string           `json:"name"`
	Email           string           `json:"email" gorm:"uniqueIndex;not null"`
	Username        string           `json:"username" gorm:"uniqueIndex;not null"`
	Password        string           `json:"password"`
	Role            string           `json:"role" gorm:"type:varchar(30);not null;default:'user'"`
	Status          string           `json:"status" gorm:"type:varchar(30);not null;default:'active'"`
	AvatarURL       string           `json:"avatar_url" gorm:"size:500"`
	Bio             string           `json:"bio"`
	Phone           string           `json:"phone" gorm:"size:30"`
	DateOfBirth     *time.Time       `json:"date_of_birth"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	// Relationships
	Channel         *Channel         `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Playlists       []Playlist       `json:"playlists" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscriptions   []Subscription   `json:"subscriptions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ViewHistories   []ViewHistory    `json:"view_histories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoReactions  []VideoReaction  `json:"video_reactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentReactions []CommentReaction `json:"comment_reactions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
