package model

import "time"

type Subscription struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null;index:idx_subscription_user_channel,unique"`
	ChannelID int       `json:"channel_id" gorm:"not null;index:idx_subscription_user_channel,unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	User      User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Channel   Channel   `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}