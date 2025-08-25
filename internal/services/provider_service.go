package services

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/utils/pagination"
	"context"
)

type ProviderServiceConfig struct {
	Repositories *repositories.Repositories
}

type ProviderService struct {
	providerRepository *repositories.ProviderRepository
}

func NewProviderService(config *ProviderServiceConfig) *ProviderService {
	return &ProviderService{
		providerRepository: config.Repositories.Provider,
	}
}

func (s *ProviderService) GetProviders(ctx context.Context, search string, with []string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.providerRepository.GetAll(ctx, with, page, limit)
	} else {
		pagination, err = s.providerRepository.Search(ctx, search, with, page, limit)
	}
	if err != nil {
		return pagination, err
	}

	providers := pagination.Data.([]models.Provider)

	// Get workspace counts
	counts, err := s.providerRepository.GetWorkspaceCounts(ctx)
	if err != nil {
		return pagination, err
	}

	pagination.Data = dto.ToProviderDTOs(providers, counts)

	return pagination, nil
}
