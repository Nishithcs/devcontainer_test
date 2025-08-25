package services

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/services/devpod"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/rabbitmq"
	"clusterix-code/internal/websocket"
	"fmt"
	"github.com/hibiken/asynq"
)

type Services struct {
	Publisher              *PublisherService
	User                   *UserService
	MachineConfig          *MachineConfigService
	Provider               *ProviderService
	GitPersonalAccessToken *GitPersonalAccessTokenService
	Repository             *RepositoryService
	Workspace              *WorkspaceService
	WorkspaceConfig        *WorkspaceConfigService
	WorkspaceLog           *WorkspaceLogService
	Socket                 *SocketService
	Devpod                 *devpod.DevpodService
}

type ServiceConfig struct {
	Repositories *repositories.Repositories
	RabbitMQ     *rabbitmq.RabbitMQ
	ApiClients   *api_clients.APIClients
	Hub          *websocket.Hub
	Redis        config.RedisConfig
}

func Provider(c *di.Container) (*Services, error) {
	repos := di.Make[*repositories.Repositories](c)
	rabbitMQ := di.Make[*rabbitmq.RabbitMQ](c)
	cfg := di.Make[*config.Config](c)
	apiClients := di.Make[*api_clients.APIClients](c)

	hub := websocket.NewHub()
	go hub.Run()

	return NewServices(&ServiceConfig{
		Repositories: repos,
		RabbitMQ:     rabbitMQ,
		ApiClients:   apiClients,
		Hub:          hub,
		Redis:        cfg.Redis,
	}), nil
}

func NewServices(config *ServiceConfig) *Services {
	publisher := NewPublisherService(&PublisherServiceConfig{
		RabbitMQ: config.RabbitMQ,
	})

	devpodService := devpod.NewDevpodService()
	fmt.Println(fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port))
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Username: config.Redis.Username,
		Password: config.Redis.Password,
	})

	socketService := NewSocketService(&SocketServiceConfig{
		Hub: config.Hub,
	})
	config.Hub.MessageRouter = socketService

	workspaceConfigService := NewWorkspaceConfigService(&WorkspaceConfigServiceConfig{
		Repositories: config.Repositories,
	})

	workspaceLogService := NewWorkspaceLogService(&WorkspaceLogServiceConfig{
		Repositories: config.Repositories,
	})

	return &Services{
		Publisher: publisher,
		User: NewUserService(&UserServiceConfig{
			Repositories: config.Repositories,
		}),
		MachineConfig: NewMachineConfigService(&MachineConfigServiceConfig{
			Repositories: config.Repositories,
		}),
		Provider: NewProviderService(&ProviderServiceConfig{
			Repositories: config.Repositories,
		}),
		GitPersonalAccessToken: NewGitPersonalAccessTokenService(&GitPersonalAccessTokenServiceConfig{
			Repositories: config.Repositories,
		}),
		Repository: NewRepositoryService(&RepositoryServiceConfig{
			Repositories: config.Repositories,
		}),
		Workspace: NewWorkspaceService(&WorkspaceServiceConfig{
			Repositories:    config.Repositories,
			Publisher:       publisher,
			Socket:          socketService,
			WorkspaceConfig: workspaceConfigService,
			Devpod:          devpodService,
			AsynqClient:     asynqClient,
		}),
		WorkspaceConfig: workspaceConfigService,
		WorkspaceLog:    workspaceLogService,
		Socket:          socketService,
		Devpod:          devpodService,
	}
}
