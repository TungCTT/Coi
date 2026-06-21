package video

import (
	"coi/internal/model"
	"coi/pkg/txmanager"
	"context"
	"time"

	"gorm.io/gorm"
)

type VideoRepository interface {
	Create(ctx context.Context, video *model.Video) error
	GetByID(ctx context.Context, id int) (*model.Video, error)
	GetByStorageKey(ctx context.Context, key string) (*model.Video, error)
	GetByChannelID(ctx context.Context, channelID int) ([]model.Video, error)
	GetByName(ctx context.Context, name string) ([]model.Video, error)
	MarkUploaded(ctx context.Context, id int, info ObjectInfo) error
	UpdateStatus(ctx context.Context, id int, status model.VideoStatus) error
	FindExpiredUploading(ctx context.Context, now time.Time) ([]model.Video, error)
	Delete(ctx context.Context, id int) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *model.Video) error {
	db := txmanager.GetTx(ctx, r.db)
	return db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) GetByID(ctx context.Context, id int) (*model.Video, error) {
	db := txmanager.GetTx(ctx, r.db)
	var video model.Video
	err := db.WithContext(ctx).First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) GetByStorageKey(ctx context.Context, key string) (*model.Video, error) {
	db := txmanager.GetTx(ctx, r.db)
	var video model.Video
	err := db.WithContext(ctx).Where("storage_key = ?", key).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) GetByChannelID(ctx context.Context, channelID int) ([]model.Video, error) {
	db := txmanager.GetTx(ctx, r.db)
	var videos []model.Video
	err := db.WithContext(ctx).Where("channel_id = ?", channelID).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (r *videoRepository) GetByName(ctx context.Context, name string) ([]model.Video, error) {
	db := txmanager.GetTx(ctx, r.db)
	var videos []model.Video
	err := db.WithContext(ctx).Where("title LIKE ?", "%"+name+"%").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (r *videoRepository) MarkUploaded(ctx context.Context, id int, info ObjectInfo) error {
	db := txmanager.GetTx(ctx, r.db)
	now := time.Now()
	return db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", id).Updates(map[string]any{
		"status":       model.VideoStatusUploaded,
		"file_size":    info.Size,
		"content_type": info.ContentType,
		"etag":         info.ETag,
		"uploaded_at":  &now,
	}).Error
}

func (r *videoRepository) UpdateStatus(ctx context.Context, id int, status model.VideoStatus) error {
	db := txmanager.GetTx(ctx, r.db)
	return db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", id).Update("status", status).Error
}

func (r *videoRepository) FindExpiredUploading(ctx context.Context, now time.Time) ([]model.Video, error) {
	db := txmanager.GetTx(ctx, r.db)
	var videos []model.Video
	err := db.WithContext(ctx).
		Where("status = ? AND upload_expires_at < ?", model.VideoStatusUploading, now).
		Find(&videos).Error
	return videos, err
}

func (r *videoRepository) Delete(ctx context.Context, id int) error {
	db := txmanager.GetTx(ctx, r.db)
	return db.WithContext(ctx).Delete(&model.Video{}, id).Error
}
