package model

import "time"

type ViewHistory struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	UserID              int       `json:"user_id" gorm:"not null;index:idx_view_history_user_video,unique"`
	VideoID             int       `json:"video_id" gorm:"not null;index:idx_view_history_user_video,unique"`
	ProgressSeconds     int64     `json:"progress_seconds" gorm:"default:0"`
	Completed           bool      `json:"completed" gorm:"default:false"`
	LastWatchedAt       time.Time `json:"last_watched_at" gorm:"autoUpdateTime"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	User                User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Video               Video     `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ViewHistory) TableName() string {
	return "view_history"
}