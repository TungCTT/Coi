package model

import "time"

type Notification struct {
	ID        int                    `json:"id" gorm:"primaryKey"`
	UserID    int                    `json:"user_id" gorm:"not null"`
	User      User                   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      NotificationType       `json:"type" gorm:"type:varchar(30);not null"`
	EntityType NotificationEntityType `json:"entity_type" gorm:"type:varchar(30);not null"`
	EntityID  int                    `json:"entity_id"`
	Message   string                 `json:"message"`
	DeepLink  string                 `json:"deep_link" gorm:"size:500"`
	IsRead    bool                   `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
}