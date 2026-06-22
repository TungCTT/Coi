package channel

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"coi/internal/model"
	"coi/internal/user"
	"coi/pkg/txmanager"

	"github.com/google/uuid"
)

const (
	channelMediaUploadTTL = 30 * time.Minute
	maxChannelAvatarSize  = 2 << 20
	maxChannelBannerSize  = 10 << 20
)

type ChannelService interface {
	CreateChannel(ctx context.Context, req *model.CreateChannelRequest) (*model.Channel, error)
	GetChannelByID(ctx context.Context, id int) (*model.Channel, error)
	GetMyChannel(ctx context.Context, userID int) (*model.Channel, error)
	GetChannelsByName(ctx context.Context, name string) ([]model.Channel, error)
	CreateMediaUploadSession(ctx context.Context, userID int, channelID int, req *model.CreateChannelMediaUploadRequest) (*model.ChannelMediaUploadResponse, error)
	ConfirmMediaUpload(ctx context.Context, userID int, channelID int, req *model.ConfirmChannelMediaUploadRequest) (*model.Channel, error)
}

type channelService struct {
	crepo     ChannelRepository
	urepo     user.UserRepository
	txManager txmanager.TxManager
	storage   ChannelMediaStorage
}

func NewChannelService(crepo ChannelRepository, urepo user.UserRepository, txManager txmanager.TxManager, storage ...ChannelMediaStorage) ChannelService {
	var mediaStorage ChannelMediaStorage
	if len(storage) > 0 {
		mediaStorage = storage[0]
	}
	return &channelService{crepo: crepo, urepo: urepo, txManager: txManager, storage: mediaStorage}
}

func (s *channelService) CreateChannel(ctx context.Context, req *model.CreateChannelRequest) (*model.Channel, error) {
	var createdChannel *model.Channel

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.crepo.ExistsByUserID(txCtx, req.UserID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("only one channel per user is allowed")
		}

		channel := &model.Channel{
			UserID:      req.UserID,
			Name:        req.Name,
			Handle:      req.Handle,
			Description: req.Description,
			AvatarURL:   req.AvatarURL,
			BannerURL:   req.BannerURL,
		}

		if err := s.crepo.Create(txCtx, channel); err != nil {
			return err
		}

		if err := s.urepo.UpdateRole(txCtx, req.UserID, model.RoleCreator); err != nil {
			return err
		}

		createdChannel = channel
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdChannel, nil
}

func (s *channelService) GetChannelByID(ctx context.Context, id int) (*model.Channel, error) {
	return s.crepo.FindByID(ctx, id)
}

func (s *channelService) GetMyChannel(ctx context.Context, userID int) (*model.Channel, error) {
	return s.crepo.FindByUserID(ctx, userID)
}

func (s *channelService) GetChannelsByName(ctx context.Context, name string) ([]model.Channel, error) {
	return s.crepo.FindByName(ctx, name)
}

func (s *channelService) CreateMediaUploadSession(ctx context.Context, userID int, channelID int, req *model.CreateChannelMediaUploadRequest) (*model.ChannelMediaUploadResponse, error) {
	if s.storage == nil {
		return nil, errors.New("media storage is not configured")
	}
	channel, err := s.ensureOwner(ctx, userID, channelID)
	if err != nil {
		return nil, err
	}
	if err := validateChannelMediaRequest(req.MediaType, req.FileSize, req.ContentType); err != nil {
		return nil, err
	}

	key := buildChannelMediaKey(channel.ID, req.MediaType, req.OriginalFileName)
	uploadURL, err := s.storage.CreatePresignedUploadURL(ctx, key, req.ContentType, channelMediaUploadTTL)
	if err != nil {
		return nil, err
	}

	return &model.ChannelMediaUploadResponse{
		UploadURL:  uploadURL,
		StorageKey: key,
		PublicURL:  s.storage.PublicURL(key),
	}, nil
}

func (s *channelService) ConfirmMediaUpload(ctx context.Context, userID int, channelID int, req *model.ConfirmChannelMediaUploadRequest) (*model.Channel, error) {
	if s.storage == nil {
		return nil, errors.New("media storage is not configured")
	}
	channel, err := s.ensureOwner(ctx, userID, channelID)
	if err != nil {
		return nil, err
	}

	expectedPrefix := fmt.Sprintf("channels/%d/%s/", channelID, req.MediaType)
	if !strings.HasPrefix(req.StorageKey, expectedPrefix) {
		return nil, errors.New("invalid storage_key for this channel")
	}

	info, err := s.storage.GetObjectInfo(ctx, req.StorageKey)
	if err != nil {
		return nil, err
	}
	if err := validateChannelMediaRequest(req.MediaType, info.Size, info.ContentType); err != nil {
		_ = s.storage.Delete(ctx, req.StorageKey)
		return nil, err
	}

	publicURL := s.storage.PublicURL(req.StorageKey)
	var oldKey string
	switch req.MediaType {
	case "avatar":
		oldKey = channel.AvatarStorageKey
		err = s.crepo.UpdateAvatar(ctx, channel.ID, publicURL, req.StorageKey)
	case "banner":
		oldKey = channel.BannerStorageKey
		err = s.crepo.UpdateBanner(ctx, channel.ID, publicURL, req.StorageKey)
	default:
		err = errors.New("invalid media_type")
	}
	if err != nil {
		return nil, err
	}
	if oldKey != "" && oldKey != req.StorageKey {
		_ = s.storage.Delete(ctx, oldKey)
	}

	return s.crepo.FindByID(ctx, channel.ID)
}

func (s *channelService) ensureOwner(ctx context.Context, userID int, channelID int) (*model.Channel, error) {
	channel, err := s.crepo.FindByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if channel.UserID != userID {
		return nil, errors.New("you do not own this channel")
	}
	return channel, nil
}

func validateChannelMediaRequest(mediaType string, fileSize int64, contentType string) error {
	if mediaType != "avatar" && mediaType != "banner" {
		return errors.New("media_type must be avatar or banner")
	}
	if fileSize <= 0 {
		return errors.New("file_size must be greater than 0")
	}

	maxSize := int64(maxChannelAvatarSize)
	if mediaType == "banner" {
		maxSize = maxChannelBannerSize
	}
	if fileSize > maxSize {
		return fmt.Errorf("%s file_size exceeds %d bytes", mediaType, maxSize)
	}

	switch strings.ToLower(contentType) {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return errors.New("content_type must be image/jpeg, image/png, or image/webp")
	}
}

func buildChannelMediaKey(channelID int, mediaType string, originalFileName string) string {
	ext := strings.ToLower(filepath.Ext(originalFileName))
	return fmt.Sprintf("channels/%d/%s/%s%s", channelID, mediaType, uuid.NewString(), sanitizeMediaExtension(ext))
}

func sanitizeMediaExtension(ext string) string {
	re := regexp.MustCompile(`[^a-z0-9.]+`)
	ext = re.ReplaceAllString(strings.ToLower(ext), "")
	if ext == "" || ext == "." {
		return ".image"
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
}
