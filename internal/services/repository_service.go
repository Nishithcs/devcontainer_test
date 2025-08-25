package services

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/utils/pagination"
	"context"
)

type RepositoryServiceConfig struct {
	Repositories *repositories.Repositories
}

type RepositoryService struct {
	gitRepository *repositories.GitRepository
}

func NewRepositoryService(config *RepositoryServiceConfig) *RepositoryService {
	return &RepositoryService{
		gitRepository: config.Repositories.GitRepository,
	}
}

func (s *RepositoryService) GetRepositories(ctx context.Context, organizationId uint32, search, status string, with []string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.gitRepository.GetRepositories(ctx, organizationId, status, with, page, limit)
	} else {
		pagination, err = s.gitRepository.Search(ctx, organizationId, search, status, with, page, limit)
	}
	if err != nil {
		return pagination, err
	}

	repos := pagination.Data.([]models.Repository)
	pagination.Data = dto.ToRepositoryDTOs(repos)

	return pagination, nil
}

func (s *RepositoryService) GetRepository(ctx context.Context, repositoryId uint64) (dto.RepositoryDTO, error) {
	repo, err := s.gitRepository.GetByID(ctx, repositoryId, []string{"User", "MachineConfig"})
	if err != nil {
		return dto.RepositoryDTO{}, err
	}
	return dto.ToRepositoryDTO(*repo), nil
}

func (s *RepositoryService) CreateRepository(ctx context.Context, req requests.CreateRepositoryRequest) (dto.RepositoryDTO, error) {
	machineConfigID := req.MachineConfigID
	repo := models.Repository{
		Title:           req.Title,
		MachineConfigID: machineConfigID,
		RepositoryURL:   req.RepositoryURL,
		CreatedByID:     req.CreatedByID,
		OrganizationID:  req.OrganizationID,
		Status:          req.Status,
		AddedByAdmin:    req.AddedByAdmin,
	}
	if err := s.gitRepository.Create(ctx, &repo); err != nil {
		return dto.RepositoryDTO{}, err
	}

	repoWithRelations, err := s.gitRepository.GetByID(ctx, repo.ID, []string{"User", "MachineConfig"})
	if err != nil {
		return dto.RepositoryDTO{}, err
	}

	return dto.ToRepositoryDTO(*repoWithRelations), nil
}

func (s *RepositoryService) UpdateRepository(ctx context.Context, req requests.UpdateRepositoryRequest) (dto.RepositoryDTO, error) {
	repo, err := s.gitRepository.GetByID(ctx, req.ID, []string{})
	if err != nil {
		return dto.RepositoryDTO{}, err
	}

	if req.Title != "" {
		repo.Title = req.Title
	}

	if req.MachineConfigID > 0 {
		repo.MachineConfigID = req.MachineConfigID
	}

	if req.RepositoryURL != "" {
		repo.RepositoryURL = req.RepositoryURL
	}
	if req.Status != "" {
		repo.Status = req.Status
	}

	if err := s.gitRepository.Update(ctx, repo); err != nil {
		return dto.RepositoryDTO{}, err
	}

	updatedRepo, err := s.gitRepository.GetByID(ctx, repo.ID, []string{"User", "MachineConfig"})
	if err != nil {
		return dto.RepositoryDTO{}, err
	}

	return dto.ToRepositoryDTO(*updatedRepo), nil
}

func (s *RepositoryService) DeleteRepository(ctx context.Context, repositoryId string) error {
	if err := s.gitRepository.DeleteRepository(ctx, repositoryId); err != nil {
		return err
	}
	return nil
}

func (s *RepositoryService) GetUserRepositories(ctx context.Context, organizationId uint32, search string, with []string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.gitRepository.GetUserRepositories(ctx, organizationId, with, page, limit)
	} else {
		pagination, err = s.gitRepository.SearchUserRepositories(ctx, organizationId, search, with, page, limit)
	}
	if err != nil {
		return pagination, err
	}

	repos := pagination.Data.([]models.Repository)
	pagination.Data = dto.ToRepositoryDTOs(repos)

	return pagination, nil
}

func (s *RepositoryService) CreateUserRepository(ctx context.Context, req requests.CreateUserRepositoryRequest) (dto.RepositoryDTO, error) {
	repo := models.Repository{
		Title:          req.Title,
		RepositoryURL:  req.RepositoryURL,
		CreatedByID:    req.CreatedByID,
		OrganizationID: req.OrganizationID,
		Status:         req.Status,
		AddedByAdmin:   req.AddedByAdmin,
	}
	if err := s.gitRepository.Create(ctx, &repo); err != nil {
		return dto.RepositoryDTO{}, err
	}
	return dto.ToRepositoryDTO(repo), nil
}
