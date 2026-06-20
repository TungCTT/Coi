package video

import (
	"context"

	"coi/internal/model"
)

type VideoRepository interface {
	Create(ctx context.Context, video *model.Video) error
	FindByID(ctx context.Context, id int) (*model.Video, error)
	FindReadyByID(ctx context.Context, id int) (*model.Video, error)
	ListReady(ctx context.Context, limit int, offset int) ([]model.Video, error)
	ListReadyByChannelID(ctx context.Context, channelID int, limit int, offset int) ([]model.Video, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateOriginalFile(ctx context.Context, id int, videoURL string, originalFileName string, fileSize int64) error
	UpdateProcessingResult(ctx context.Context, id int, videoURL string, thumbnailURL string, durationSeconds int64, status string) error
	IncrementViewCount(ctx context.Context, id int) error
}
