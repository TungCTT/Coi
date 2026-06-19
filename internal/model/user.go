package model

import (
	"time"

	"coi/pkg/jwt"
)

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

// ─── Request DTOs ───────────────────────────────────────────────────────────

// CreateUserRequest là dữ liệu client gửi lên khi đăng ký.
// Tag `binding:"required"` là của Gin validator — tự động trả 400 nếu thiếu field.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest là dữ liệu client gửi lên khi đăng nhập.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ─── Response DTOs ───────────────────────────────────────────────────────────

// UserResponse là thông tin user trả về cho client.
// Không bao gồm Password — KHÔNG BAO GIỜ trả password (dù đã hash) về client.
type UserResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	AvatarURL string     `json:"avatar_url"`
	Bio       string     `json:"bio"`
	CreatedAt time.Time  `json:"created_at"`
}

// AuthResponse gộp token pair + thông tin user vào một response duy nhất.
// Client dùng AccessToken để gọi các API bảo vệ, RefreshToken để lấy token mới.
type AuthResponse struct {
	*jwt.TokenPair
	User UserResponse `json:"user"`
}
