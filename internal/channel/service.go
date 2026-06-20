package channel

import (
	"context"
	"errors"

	"coi/internal/model"
	"coi/internal/user"
	"coi/pkg/txmanager"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, req *model.CreateChannelRequest) (*model.Channel, error)
	GetChannelByID(ctx context.Context, id int) (*model.Channel, error)
	GetChannelsByName(ctx context.Context, name string) ([]model.Channel, error)
}

type channelService struct {
	crepo     ChannelRepository
	urepo     user.UserRepository
	txManager txmanager.TxManager
}

func NewChannelService(crepo ChannelRepository, urepo user.UserRepository, txManager txmanager.TxManager) ChannelService {
	return &channelService{crepo: crepo, urepo: urepo, txManager: txManager}
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

func (s *channelService) GetChannelsByName(ctx context.Context, name string) ([]model.Channel, error) {
	return s.crepo.FindByName(ctx, name)
}
