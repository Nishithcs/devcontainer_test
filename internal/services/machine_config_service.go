package services

import (
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/utils/pagination"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type MachineConfigServiceConfig struct {
	Repositories *repositories.Repositories
}

type MachineConfigService struct {
	machineConfigRepository *repositories.MachineConfigRepository
}

func NewMachineConfigService(config *MachineConfigServiceConfig) *MachineConfigService {
	return &MachineConfigService{
		machineConfigRepository: config.Repositories.MachineConfig,
	}
}

func (s *MachineConfigService) Import(ctx context.Context) error {
	jsonPath := filepath.Join("internal", "data", "static", "machine_configs.json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read machine configs file: %w", err)
	}

	var configs []models.MachineConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to unmarshal machine configs: %w", err)
	}

	for _, config := range configs {
		fmt.Printf("%+v\n", config)
		existing, err := s.machineConfigRepository.GetByInstanceType(ctx, config.InstanceType)
		if err != nil {
			return fmt.Errorf("failed to check existing config: %w", err)
		}

		if existing == nil {
			if err := s.machineConfigRepository.CreateConfig(ctx, &config); err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}
		} else {
			config.ID = existing.ID
			if err := s.machineConfigRepository.UpdateConfig(ctx, &config); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}
		}
	}

	fmt.Println("Successfully imported machine configs")
	return nil
}

func (s *MachineConfigService) GetMachineConfigs(ctx context.Context, search string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.machineConfigRepository.GetAll(ctx, page, limit)
	} else {
		pagination, err = s.machineConfigRepository.Search(ctx, search, page, limit)
	}
	if err != nil {
		pagination = pagination
	}

	machines := pagination.Data.([]models.MachineConfig)
	pagination.Data = dto.ToMachineConfigDTOs(machines)

	return pagination, nil
}

func (s *MachineConfigService) GetMachineConfig(ctx context.Context, userId uint64, machineConfigId uint64) (dto.MachineConfigDTO, error) {
	machineConfig, err := s.machineConfigRepository.GetByID(ctx, machineConfigId)
	if err != nil {
		return dto.MachineConfigDTO{}, err
	}
	return dto.ToMachineConfigDTO(*machineConfig), nil
}
