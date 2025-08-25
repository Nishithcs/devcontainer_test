package api_clients

import (
	"clusterix-code/internal/auth"
	"clusterix-code/internal/config"
	"clusterix-code/internal/utils/di"
)

type APIClients struct {
	Auth *AuthAPIClient
}

type Config struct {
	Auth GeneralAPIConfig
}

type GeneralAPIConfig struct {
	BaseURL     string
	MasterToken string
}

func New(config Config) (*APIClients, error) {
	masterTokenProvider := auth.NewStaticTokenProvider(config.Auth.MasterToken)

	return &APIClients{
		Auth: NewAuthAPIClient(
			config.Auth.BaseURL,
			masterTokenProvider,
		),
	}, nil
}

func Provider(c *di.Container) (*APIClients, error) {
	cfg := di.Make[*config.Config](c)
	return New(Config{
		Auth: GeneralAPIConfig{
			BaseURL:     cfg.ExternalServices.Auth.BaseURL,
			MasterToken: cfg.ExternalServices.Auth.MasterToken,
		},
	})
}
