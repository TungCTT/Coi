package video

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"coi/internal/category"
	"coi/internal/channel"
	"coi/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
)

const (
	maxVideoUploadSize = 5 << 30
	uploadURLTTL       = 2 * time.Hour
)

type VideoService interface {
	CreateUploadSession(ctx context.Context, userID int, req *model.CreateVideoUploadRequest) (*model.CreateVideoUploadResponse, error)
	ConfirmUpload(ctx context.Context, userID int, videoID int) (*model.Video, error)
	HandleUploadCompleted(ctx context.Context, storageKey string) (*model.Video, error)
	OpenPublicVideo(ctx context.Context, id int, byteRange string) (*model.Video, io.ReadCloser, *ObjectInfo, error)
	ListPublicVideos(ctx context.Context, query string, limit int, offset int) ([]model.Video, error)
	ListPublicVideosByChannel(ctx context.Context, channelID int, limit int, offset int) ([]model.Video, error)
	GetPublicVideo(ctx context.Context, id int) (*model.Video, error)
	DeleteVideo(ctx context.Context, userID int, videoID int) error
	ReconcileExpiredUploads(ctx context.Context, now time.Time) error
}

type videoService struct {
	videoRepo    VideoRepository
	storage      VideoStorage
	channelRepo  channel.ChannelRepository
	categoryRepo category.CategoryRepository
}

func NewVideoService(
	videoRepo VideoRepository,
	storage VideoStorage,
	channelRepo channel.ChannelRepository,
	categoryRepo category.CategoryRepository,
) VideoService {
	return &videoService{
		videoRepo:    videoRepo,
		storage:      storage,
		channelRepo:  channelRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *videoService) CreateUploadSession(ctx context.Context, userID int, req *model.CreateVideoUploadRequest) (*model.CreateVideoUploadResponse, error) {
	if err := validateUploadRequest(req); err != nil {
		return nil, err
	}

	if err := s.ensureChannelOwner(ctx, req.ChannelID, userID); err != nil {
		return nil, err
	}

	if _, err := s.categoryRepo.FindByID(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	key := buildVideoStorageKey(req.ChannelID, req.OriginalFileName)
	now := time.Now()
	video := &model.Video{
		ChannelID:        req.ChannelID,
		CategoryID:       req.CategoryID,
		Title:            strings.TrimSpace(req.Title),
		Description:      strings.TrimSpace(req.Description),
		Visibility:       req.Visibility,
		StorageKey:       key,
		VideoURL:         s.storage.PublicURL(key),
		OriginalFileName: req.OriginalFileName,
		FileSize:         req.FileSize,
		ContentType:      req.ContentType,
		Status:           model.VideoStatusUploading,
		UploadExpiresAt:  now.Add(uploadURLTTL),
		CommentsEnabled:  true,
	}

	uploadURL, err := s.storage.CreatePresignedUploadURL(ctx, key, req.ContentType, uploadURLTTL)
	if err != nil {
		return nil, err
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	return &model.CreateVideoUploadResponse{
		VideoID:         video.ID,
		UploadURL:       uploadURL,
		StorageKey:      key,
		Status:          string(video.Status),
		UploadExpiresAt: video.UploadExpiresAt,
	}, nil
}

func (s *videoService) ConfirmUpload(ctx context.Context, userID int, videoID int) (*model.Video, error) {
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureChannelOwner(ctx, video.ChannelID, userID); err != nil {
		return nil, err
	}

	return s.markUploadedFromStorage(ctx, video)
}

func (s *videoService) HandleUploadCompleted(ctx context.Context, storageKey string) (*model.Video, error) {
	if strings.TrimSpace(storageKey) == "" {
		return nil, errors.New("storage_key is required")
	}

	video, err := s.videoRepo.GetByStorageKey(ctx, storageKey)
	if err != nil {
		return nil, err
	}

	return s.markUploadedFromStorage(ctx, video)
}

func (s *videoService) OpenPublicVideo(ctx context.Context, id int, byteRange string) (*model.Video, io.ReadCloser, *ObjectInfo, error) {
	video, err := s.videoRepo.GetPublicReadyByID(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	out, err := s.storage.Open(ctx, video.StorageKey, byteRange)
	if err != nil {
		return nil, nil, nil, err
	}

	info := &ObjectInfo{
		Key:          video.StorageKey,
		Size:         aws.ToInt64(out.ContentLength),
		ContentType:  strings.TrimSpace(aws.ToString(out.ContentType)),
		ETag:         strings.Trim(aws.ToString(out.ETag), "\""),
		ContentRange: aws.ToString(out.ContentRange),
	}
	if info.ContentType == "" {
		info.ContentType = video.ContentType
	}

	return video, out.Body, info, nil
}

func (s *videoService) ListPublicVideos(ctx context.Context, query string, limit int, offset int) ([]model.Video, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.videoRepo.ListPublicReady(ctx, strings.TrimSpace(query), limit, offset)
}

func (s *videoService) ListPublicVideosByChannel(ctx context.Context, channelID int, limit int, offset int) ([]model.Video, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if _, err := s.channelRepo.FindByID(ctx, channelID); err != nil {
		return nil, err
	}
	return s.videoRepo.ListPublicReadyByChannelID(ctx, channelID, limit, offset)
}

func (s *videoService) GetPublicVideo(ctx context.Context, id int) (*model.Video, error) {
	return s.videoRepo.GetPublicReadyByID(ctx, id)
}

func (s *videoService) DeleteVideo(ctx context.Context, userID int, videoID int) error {
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return err
	}

	if err := s.ensureChannelOwner(ctx, video.ChannelID, userID); err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, video.StorageKey); err != nil {
		return err
	}

	return s.videoRepo.Delete(ctx, video.ID)
}

func (s *videoService) ReconcileExpiredUploads(ctx context.Context, now time.Time) error {
	videos, err := s.videoRepo.FindExpiredUploading(ctx, now)
	if err != nil {
		return err
	}

	for _, v := range videos {
		exists, err := s.storage.Exists(ctx, v.StorageKey)
		if err != nil {
			continue
		}
		if exists {
			_, _ = s.markUploadedFromStorage(ctx, &v)
			continue
		}
		_ = s.videoRepo.UpdateStatus(ctx, v.ID, model.VideoStatusFailed)
	}

	return nil
}

func (s *videoService) markUploadedFromStorage(ctx context.Context, video *model.Video) (*model.Video, error) {
	if video.Status == model.VideoStatusReady {
		return video, nil
	}
	if video.Status == model.VideoStatusUploaded {
		if err := s.videoRepo.UpdateStatus(ctx, video.ID, model.VideoStatusReady); err != nil {
			return nil, err
		}
		return s.videoRepo.GetByID(ctx, video.ID)
	}
	if video.Status != model.VideoStatusUploading {
		return video, nil
	}

	info, err := s.storage.GetObjectInfo(ctx, video.StorageKey)
	if err != nil {
		return nil, err
	}

	if err := s.videoRepo.MarkUploaded(ctx, video.ID, *info); err != nil {
		return nil, err
	}

	return s.videoRepo.GetByID(ctx, video.ID)
}

func (s *videoService) ensureChannelOwner(ctx context.Context, channelID int, userID int) error {
	ch, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return err
	}
	if ch.UserID != userID {
		return errors.New("you do not own this channel")
	}
	return nil
}

func validateUploadRequest(req *model.CreateVideoUploadRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}
	if req.FileSize <= 0 {
		return errors.New("file_size must be greater than 0")
	}
	if req.FileSize > maxVideoUploadSize {
		return fmt.Errorf("file_size exceeds %d bytes", maxVideoUploadSize)
	}
	if !strings.HasPrefix(strings.ToLower(req.ContentType), "video/") {
		return errors.New("content_type must be a video type")
	}
	if strings.TrimSpace(req.OriginalFileName) == "" {
		return errors.New("original_file_name is required")
	}
	return nil
}

func buildVideoStorageKey(channelID int, originalFileName string) string {
	fileName := sanitizeFileName(originalFileName)
	return fmt.Sprintf("videos/%d/%s-%s", channelID, uuid.NewString(), fileName)
}

func sanitizeFileName(name string) string {
	base := filepath.Base(name)
	base = strings.TrimSpace(base)
	if base == "." || base == "" {
		return "video"
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	base = re.ReplaceAllString(base, "-")
	base = strings.Trim(base, ".-")
	if base == "" {
		return "video"
	}
	return base
}
