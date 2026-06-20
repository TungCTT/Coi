package channel

import (
	"context"

	"coi/internal/model"
	"coi/pkg/txmanager"

	"gorm.io/gorm"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *model.Channel) error
	FindByID(ctx context.Context, id int) (*model.Channel, error)
	FindByUserID(ctx context.Context, userID int) (*model.Channel, error)
	FindByName(ctx context.Context, name string) ([]model.Channel, error)
	ExistsByUserID(ctx context.Context, userID int) (bool, error)
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(ctx context.Context, channel *model.Channel) error {
	db := txmanager.GetTx(ctx, r.db)
	return db.WithContext(ctx).Create(channel).Error
}

func (r *channelRepository) FindByID(ctx context.Context, id int) (*model.Channel, error) {
	db := txmanager.GetTx(ctx, r.db)
	var channel model.Channel
	err := db.WithContext(ctx).First(&channel, id).Error
	if err != nil {
		return nil, err
	}
	return &channel, err
}

func (r *channelRepository) FindByUserID(ctx context.Context, userID int) (*model.Channel, error) {
	db := txmanager.GetTx(ctx, r.db)
	var channel model.Channel
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *channelRepository) FindByName(ctx context.Context, name string) ([]model.Channel, error) {
	db := txmanager.GetTx(ctx, r.db)
	var channels []model.Channel
	err := db.WithContext(ctx).Where("name = ?", name).Order("subscribers DESC").Find(&channels).Error
	return channels, err
}
func (r *channelRepository) ExistsByUserID(ctx context.Context, userID int) (bool, error) {
	db := txmanager.GetTx(ctx, r.db)
	var count int64
	err := db.WithContext(ctx).Model(&model.Channel{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
