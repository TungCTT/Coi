package user

import (
	"errors"
	"coi/internal/model"
	jwtpkg "coi/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *model.CreateUserRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(refreshToken string) (*jwtpkg.TokenPair, error)
	GetUserByID(id int) (*model.UserResponse, error)
}


type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

// Register tạo tài khoản mới
func (s *authService) Register(req *model.CreateUserRequest) (*model.AuthResponse, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email đã được sử dụng")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Lỗi DB thực sự (kết nối, v.v.) — không phải "không tìm thấy"
		return nil, err
	}

	_, err = s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username đã được sử dụng")
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
		Password: string(hashedPassword), // Chỉ lưu hash, không lưu plain-text
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
			return nil, errors.New("email hoặc mật khẩu không đúng")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}

	return buildAuthResponse(user)
}

// RefreshToken nhận refresh token, validate và cấp access token mới.
func (s *authService) RefreshToken(refreshToken string) (*jwtpkg.TokenPair, error) {
	claims, err := jwtpkg.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token không hợp lệ hoặc đã hết hạn")
	}


	tokenPair, err := jwtpkg.GenerateTokenPair(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

// GetUserByID lấy thông tin user theo ID (dùng cho route /me).
func (s *authService) GetUserByID(id int) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("không tìm thấy user")
		}
		return nil, err
	}
	resp := toUserResponse(user)
	return &resp, nil
}


// buildAuthResponse là helper nội bộ tạo AuthResponse từ User entity.
// Dùng chung cho cả Register và Login để tránh lặp code.
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

// toUserResponse chuyển đổi User entity (có password hash) thành UserResponse (không có password).
// Đây là bước quan trọng để đảm bảo không bao giờ lộ password ra ngoài.
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
