package services

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/utils/logger"
	"context"
	"go.uber.org/zap"
	"strings"
)

type UserServiceConfig struct {
	Repositories *repositories.Repositories
	ApiClients   *api_clients.APIClients
}

type UserService struct {
	userRepository *repositories.UserRepository
	apiClients     *api_clients.APIClients
}

func NewUserService(config *UserServiceConfig) *UserService {
	return &UserService{
		userRepository: config.Repositories.User,
		apiClients:     config.ApiClients,
	}
}

func (s *UserService) SyncUser(ctx context.Context, user *dto.AuthUserDto) (*models.User, error) {
	if user.Profile == nil {
		user.Profile = &dto.AuthUserProfile{}
		names := strings.Split(user.Name, " ")
		if len(names) > 0 {
			user.Profile.FirstName = names[0]
			if len(names) > 1 {
				user.Profile.LastName = strings.Join(names[1:], " ")
			}
		}
	}
	userModel, err := s.userRepository.SyncUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (s *UserService) SyncCreatedUser(ctx context.Context, userID uint64) error {
	authUserDto, err := s.apiClients.Auth.GetUser(ctx, userID)
	if err != nil {
		logger.Error("error in getting user from auth service", err, zap.Uint64("user_id", userID))
		return err
	}

	_, err = s.SyncUser(ctx, authUserDto)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) SyncUpdatedUser(ctx context.Context, userID uint64) error {
	authUserDto, err := s.apiClients.Auth.GetUser(ctx, userID)
	if err != nil {
		logger.Error("error in getting user from auth service", err, zap.Uint64("user_id", userID))
		return err
	}

	_, err = s.SyncUser(ctx, authUserDto)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) SyncDeletedUser(ctx context.Context, userID uint64) error {
	err := s.userRepository.SoftDelete(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
