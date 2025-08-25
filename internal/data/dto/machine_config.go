package dto

import "clusterix-code/internal/data/models"

type MachineConfigDTO struct {
	ID                 uint64  `json:"id"`
	InstanceType       string  `json:"instance_type"`
	Category           string  `json:"category"`
	CPUCores           int     `json:"cpu_cores"`
	MemoryGB           float64 `json:"memory_gb"`
	StorageType        string  `json:"storage_type"`
	StorageSizeGB      int     `json:"storage_size_gb"`
	NetworkPerformance string  `json:"network_performance"`
	Architecture       string  `json:"architecture"`
	Hypervisor         string  `json:"hypervisor"`
	Generation         string  `json:"generation"`
	EnhancedNetworking bool    `json:"enhanced_networking"`
	GPU                string  `json:"gpu"`
	AdditionalFeatures string  `json:"additional_features"`
}

func ToMachineConfigDTO(config models.MachineConfig) MachineConfigDTO {
	return MachineConfigDTO{
		ID:                 config.ID,
		InstanceType:       config.InstanceType,
		Category:           config.Category,
		CPUCores:           config.CPUCores,
		MemoryGB:           config.MemoryGB,
		StorageType:        config.StorageType,
		StorageSizeGB:      config.StorageSizeGB,
		NetworkPerformance: config.NetworkPerformance,
		Architecture:       config.Architecture,
		Hypervisor:         config.Hypervisor,
		Generation:         config.Generation,
		EnhancedNetworking: config.EnhancedNetworking,
		GPU:                config.GPU,
		AdditionalFeatures: config.AdditionalFeatures,
	}
}

func ToMachineConfigDTOs(configs []models.MachineConfig) []MachineConfigDTO {
	result := make([]MachineConfigDTO, len(configs))
	for i, config := range configs {
		result[i] = ToMachineConfigDTO(config)
	}
	return result
}
