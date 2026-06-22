package user

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"coi/internal/model"
	jwtpkg "coi/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userMediaUploadTTL = 30 * time.Minute
	maxAvatarSize      = 2 << 20
)

type AuthService interface {
	Register(req *model.CreateUserRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(refreshToken string) (*jwtpkg.TokenPair, error)
	GetUserByID(id int) (*model.UserResponse, error)
	CreateMediaUploadSession(ctx context.Context, userID int, req *model.CreateUserMediaUploadRequest) (*model.UserMediaUploadResponse, error)
	ConfirmMediaUpload(ctx context.Context, userID int, req *model.ConfirmUserMediaUploadRequest) (*model.UserResponse, error)
}

type authService struct {
	repo    UserRepository
	storage UserMediaStorage
}

func NewAuthService(repo UserRepository, storage ...UserMediaStorage) AuthService {
	var mediaStorage UserMediaStorage
	if len(storage) > 0 {
		mediaStorage = storage[0]
	}
	return &authService{repo: repo, storage: mediaStorage}
}

func (s *authService) Register(req *model.CreateUserRequest) (*model.AuthResponse, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email da duoc su dung")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_, err = s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username da duoc su dung")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return buildAuthResponse(user)
}

func (s *authService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email hoac mat khau khong dung")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email hoac mat khau khong dung")
	}

	return buildAuthResponse(user)
}

func (s *authService) RefreshToken(refreshToken string) (*jwtpkg.TokenPair, error) {
	claims, err := jwtpkg.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token khong hop le hoac da het han")
	}

	tokenPair, err := jwtpkg.GenerateTokenPair(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

func (s *authService) GetUserByID(id int) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("khong tim thay user")
		}
		return nil, err
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *authService) CreateMediaUploadSession(ctx context.Context, userID int, req *model.CreateUserMediaUploadRequest) (*model.UserMediaUploadResponse, error) {
	if s.storage == nil {
		return nil, errors.New("media storage is not configured")
	}
	if err := validateUserAvatarRequest(req.FileSize, req.ContentType); err != nil {
		return nil, err
	}
	if _, err := s.repo.FindByID(userID); err != nil {
		return nil, err
	}

	key := buildUserAvatarKey(userID, req.OriginalFileName)
	uploadURL, err := s.storage.CreatePresignedUploadURL(ctx, key, req.ContentType, userMediaUploadTTL)
	if err != nil {
		return nil, err
	}

	return &model.UserMediaUploadResponse{
		UploadURL:  uploadURL,
		StorageKey: key,
		PublicURL:  s.storage.PublicURL(key),
	}, nil
}

func (s *authService) ConfirmMediaUpload(ctx context.Context, userID int, req *model.ConfirmUserMediaUploadRequest) (*model.UserResponse, error) {
	if s.storage == nil {
		return nil, errors.New("media storage is not configured")
	}

	expectedPrefix := fmt.Sprintf("users/%d/avatar/", userID)
	if !strings.HasPrefix(req.StorageKey, expectedPrefix) {
		return nil, errors.New("invalid storage_key for this user")
	}

	info, err := s.storage.GetObjectInfo(ctx, req.StorageKey)
	if err != nil {
		return nil, err
	}
	if err := validateUserAvatarRequest(info.Size, info.ContentType); err != nil {
		_ = s.storage.Delete(ctx, req.StorageKey)
		return nil, err
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	publicURL := s.storage.PublicURL(req.StorageKey)
	oldKey := user.AvatarStorageKey
	err = s.repo.UpdateAvatar(ctx, userID, publicURL, req.StorageKey)
	if err != nil {
		return nil, err
	}

	if oldKey != "" && oldKey != req.StorageKey {
		_ = s.storage.Delete(ctx, oldKey)
	}

	updated, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	resp := toUserResponse(updated)
	return &resp, nil
}

func buildAuthResponse(user *model.User) (*model.AuthResponse, error) {
	tokenPair, err := jwtpkg.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{
		TokenPair: tokenPair,
		User:      toUserResponse(user),
	}, nil
}

func toUserResponse(u *model.User) model.UserResponse {
	return model.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Username:  u.Username,
		Role:      u.Role,
		Status:    u.Status,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
	}
}

func validateUserAvatarRequest(fileSize int64, contentType string) error {
	if fileSize <= 0 {
		return errors.New("file_size must be greater than 0")
	}
	if fileSize > maxAvatarSize {
		return fmt.Errorf("avatar file_size exceeds %d bytes", maxAvatarSize)
	}

	switch strings.ToLower(contentType) {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return errors.New("content_type must be image/jpeg, image/png, or image/webp")
	}
}

func buildUserAvatarKey(userID int, originalFileName string) string {
	ext := strings.ToLower(filepath.Ext(originalFileName))
	return fmt.Sprintf("users/%d/avatar/%s%s", userID, uuid.NewString(), sanitizeMediaExtension(ext))
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
