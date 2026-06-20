package user

import (
	"context"

	"coi/internal/model"
	"coi/pkg/txmanager"

	"gorm.io/gorm"
)

// ─── Interface ───────────────────────────────────────────────────────────────
//
// UserRepository định nghĩa "hợp đồng" (contract) giữa service và database.
//
// Lý do dùng interface thay vì struct trực tiếp:
//  1. Testability: Có thể mock repository trong unit test mà không cần DB thật.
//  2. Dependency Inversion (SOLID): Service phụ thuộc vào abstraction, không phụ
//     thuộc vào GORM cụ thể. Sau này muốn đổi sang MongoDB chỉ cần viết impl mới.
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	UpdateRole(ctx context.Context, userID int, role model.UserRole) error
}

// ─── Implementation ──────────────────────────────────────────────────────────

// userRepository là concrete implementation của UserRepository dùng GORM.
// Lưu ý: struct này unexported (chữ thường) — bên ngoài package chỉ thấy interface.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository là constructor — pattern phổ biến trong Go để tạo struct
// và trả về interface (ẩn implementation details với bên ngoài).
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create lưu user mới vào DB. GORM tự điền CreatedAt, UpdatedAt.
func (r *userRepository) Create(user *model.User) error {
	// db.Create() nhận pointer, tự map fields vào câu INSERT SQL.
	return r.db.Create(user).Error
}

// FindByEmail tìm user theo email. Trả về gorm.ErrRecordNotFound nếu không có.
// Service sẽ dùng lỗi này để biết email chưa tồn tại (đăng ký) hay không tìm thấy (đăng nhập).
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	// db.Where().First() tương đương: SELECT * FROM users WHERE email = ? LIMIT 1
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername tìm user theo username. Dùng để kiểm tra username đã tồn tại chưa khi đăng ký.
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID tìm user theo primary key. Dùng cho route /me sau khi decode JWT.
func (r *userRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	// db.First(&model, id) tương đương: SELECT * FROM users WHERE id = ? LIMIT 1
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int, role model.UserRole) error {
	db := txmanager.GetTx(ctx, r.db)
	result := db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
