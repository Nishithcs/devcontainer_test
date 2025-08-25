package services

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/utils/pagination"
	"context"
)

type GitPersonalAccessTokenServiceConfig struct {
	Repositories *repositories.Repositories
}

type GitPersonalAccessTokenService struct {
	gitPersonalAccessTokenRepository *repositories.GitPersonalAccessTokenRepository
}

func NewGitPersonalAccessTokenService(config *GitPersonalAccessTokenServiceConfig) *GitPersonalAccessTokenService {
	return &GitPersonalAccessTokenService{
		gitPersonalAccessTokenRepository: config.Repositories.GitPersonalAccessToken,
	}
}

func (s *GitPersonalAccessTokenService) GetUserAccessTokens(ctx context.Context, userId uint64, search string, with []string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.gitPersonalAccessTokenRepository.GetUserAccessTokens(ctx, userId, with, page, limit)
	} else {
		pagination, err = s.gitPersonalAccessTokenRepository.Search(ctx, userId, search, with, page, limit)
	}
	if err != nil {
		return pagination, err
	}

	tokens := pagination.Data.([]models.GitPersonalAccessToken)
	pagination.Data = dto.ToGitAccessTokenDTOs(tokens)

	return pagination, nil
}

func (s *GitPersonalAccessTokenService) GetUserAccessToken(ctx context.Context, userId uint64, accessTokeId uint64) (dto.GitAccessTokenDTO, error) {
	accessToken, err := s.gitPersonalAccessTokenRepository.GetByID(ctx, userId, accessTokeId)
	if err != nil {
		return dto.GitAccessTokenDTO{}, err
	}
	return dto.ToGitAccessTokenDTO(*accessToken), nil
}

func (s *GitPersonalAccessTokenService) CreateUserAccessToken(ctx context.Context, req requests.CreateGitAccessTokenRequest) (dto.GitAccessTokenDTO, error) {
	gitAccessToken := models.GitPersonalAccessToken{
		Title:     req.Title,
		Token:     req.Token,
		UserID:    req.UserID,
		IsDefault: *req.IsDefault,
	}
	if err := s.gitPersonalAccessTokenRepository.Create(ctx, &gitAccessToken); err != nil {
		return dto.GitAccessTokenDTO{}, err
	}
	return dto.ToGitAccessTokenDTO(gitAccessToken), nil
}

func (s *GitPersonalAccessTokenService) UpdateUserAccessToken(ctx context.Context, userId uint64, req requests.UpdateGitAccessTokenRequest) (dto.GitAccessTokenDTO, error) {
	repo, err := s.gitPersonalAccessTokenRepository.GetByID(ctx, userId, req.ID)
	if err != nil {
		return dto.GitAccessTokenDTO{}, err
	}

	if req.Title != "" {
		repo.Title = req.Title
	}
	if req.IsDefault != nil {
		repo.IsDefault = *req.IsDefault
	}

	if err := s.gitPersonalAccessTokenRepository.Update(ctx, repo); err != nil {
		return dto.GitAccessTokenDTO{}, err
	}

	updatedToken, err := s.gitPersonalAccessTokenRepository.GetByID(ctx, userId, repo.ID)
	if err != nil {
		return dto.GitAccessTokenDTO{}, err
	}

	return dto.ToGitAccessTokenDTO(*updatedToken), nil
}

func (s *GitPersonalAccessTokenService) DeleteUserAccessToken(ctx context.Context, accessTokenId string) error {
	if err := s.gitPersonalAccessTokenRepository.DeleteUserAccessTokens(ctx, accessTokenId); err != nil {
		return err
	}
	return nil
}
