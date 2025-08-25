package services

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/dto"
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/data/models"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/services/devpod"
	"clusterix-code/internal/tasks"
	"clusterix-code/internal/utils/aws"
	"clusterix-code/internal/utils/pagination"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"strconv"
	"time"
)

type WorkspaceServiceConfig struct {
	Repositories    *repositories.Repositories
	Publisher       *PublisherService
	Socket          *SocketService
	WorkspaceConfig *WorkspaceConfigService
	Devpod          *devpod.DevpodService
	AsynqClient     *asynq.Client
}

type WorkspaceService struct {
	workspaceRepository            *repositories.WorkspaceRepository
	publisherService               *PublisherService
	socketService                  *SocketService
	workspaceConfigService         *WorkspaceConfigService
	devpod                         *devpod.DevpodService
	workspaceStatusEventRepository *repositories.WorkspaceStatusEventRepository
	asynqClient                    *asynq.Client
}

func NewWorkspaceService(config *WorkspaceServiceConfig) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository:            config.Repositories.Workspace,
		publisherService:               config.Publisher,
		socketService:                  config.Socket,
		workspaceConfigService:         config.WorkspaceConfig,
		devpod:                         config.Devpod,
		workspaceStatusEventRepository: config.Repositories.WorkspaceStatusEvent,
		asynqClient:                    config.AsynqClient,
	}
}

func (s *WorkspaceService) GetWorkspaces(ctx context.Context, userId uint64, search string, with []string, page, limit int) (pagination.Pagination, error) {
	var pagination pagination.Pagination
	var err error

	if search == "" {
		pagination, err = s.workspaceRepository.GetWorkspaces(ctx, userId, with, page, limit)
	} else {
		pagination, err = s.workspaceRepository.Search(ctx, userId, search, with, page, limit)
	}
	if err != nil {
		return pagination, err
	}

	workspaces := pagination.Data.([]models.Workspace)
	pagination.Data = dto.ToWorkspaceDTOs(workspaces)

	return pagination, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId uint64) (dto.WorkspaceDTO, error) {
	workspace, err := s.workspaceRepository.GetByID(ctx, workspaceId)
	if err != nil {
		return dto.WorkspaceDTO{}, err
	}
	return dto.ToWorkspaceDTO(*workspace), nil
}

func (s *WorkspaceService) GetWorkspaceIncludingDeleted(ctx context.Context, id uint64) (dto.WorkspaceDTO, error) {
	workspace, err := s.workspaceRepository.GetByIDIncludingDeleted(ctx, id)
	if err != nil {
		return dto.WorkspaceDTO{}, err
	}
	return dto.ToWorkspaceDTO(*workspace), nil
}

func (s *WorkspaceService) GetWorkspaceByFingerprint(ctx context.Context, workspaceFingerprint string) (dto.WorkspaceDTO, error) {
	workspace, err := s.workspaceRepository.GetByFingerprint(ctx, workspaceFingerprint)
	if err != nil {
		return dto.WorkspaceDTO{}, err
	}
	return dto.ToWorkspaceDTO(*workspace), nil
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, req requests.CreateWorkspaceRequest) (dto.WorkspaceDTO, error) {
	fingerprint := s.GenerateFingerprint(req.Title, req.UserID, req.OrganizationID)

	var workspaceConfigRequest requests.CreateWorkspaceConfigRequest
	workspaceConfig, err := s.workspaceConfigService.CreateWorkspaceConfig(ctx, workspaceConfigRequest)
	if err != nil {
		return dto.WorkspaceDTO{}, err
	}

	workspace := models.Workspace{
		Title:                    req.Title,
		Color:                    req.Color,
		Ide:                      req.IDE,
		RepositoryID:             req.RepositoryID,
		UserID:                   req.UserID,
		GitPersonalAccessTokenID: req.GitAccessTokenID,
		OrganizationID:           req.OrganizationID,
		Status:                   enums.WorkspaceStatus(req.Status),
		Tags:                     req.Tags,
		URL:                      "",
		ProviderID:               &req.ProviderID,
		Fingerprint:              fingerprint,
		WorkspaceConfigID:        &workspaceConfig.ID,
	}

	if err := s.workspaceRepository.Create(ctx, &workspace); err != nil {
		return dto.WorkspaceDTO{}, err
	}

	err = s.UpdateWorkspaceStatus(ctx, workspace.ID, enums.WorkspaceStatusPending, "Workspace creation is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	createdWorkspace, _ := s.workspaceRepository.GetByID(ctx, workspace.ID)

	task, err := tasks.NewStartWorkspaceTask(workspace.ID, req.UserID)
	if err != nil {
		return dto.WorkspaceDTO{}, fmt.Errorf("failed to create workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return dto.WorkspaceDTO{}, fmt.Errorf("failed to enqueue workspace creation task: %w", err)
	}
	log.Printf("✅ Enqueued workspace creation task: ID=%s queue=%s", info.ID, info.Queue)

	// Socket message
	message := dto.Message{
		EventType: constants.WorkspaceCreated,
		Channel:   fmt.Sprintf("workspace_%d_logs", workspace.ID),
		Data:      dto.ToWorkspaceDTO(*createdWorkspace),
	}
	s.socketService.SendMessage(message)

	return dto.ToWorkspaceDTO(*createdWorkspace), nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, req requests.UpdateWorkspaceRequest) (dto.WorkspaceDTO, error) {
	workspace, err := s.workspaceRepository.GetByID(ctx, req.ID)
	if err != nil {
		return dto.WorkspaceDTO{}, err
	}

	if req.Title != "" {
		workspace.Title = req.Title
	}
	if req.Color != "" {
		workspace.Color = req.Color
	}
	if req.IDE != "" {
		workspace.Ide = req.IDE
	}
	if req.RepositoryID > 0 {
		workspace.RepositoryID = req.RepositoryID
	}
	if req.GitAccessTokenID > 0 {
		workspace.GitPersonalAccessTokenID = req.GitAccessTokenID
	}
	if req.ProviderID != 0 {
		workspace.ProviderID = &req.ProviderID
	}

	if err := s.workspaceRepository.Update(ctx, workspace); err != nil {
		return dto.WorkspaceDTO{}, err
	}

	task, err := tasks.NewRebuildWorkspaceTask(req.ID, workspace.UserID)
	if err != nil {
		return dto.WorkspaceDTO{}, fmt.Errorf("failed to rebuild workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return dto.WorkspaceDTO{}, fmt.Errorf("failed to enqueue workspace rebuilding task: %w", err)
	}
	log.Printf("✅ Enqueued workspace rebuilding task: ID=%s queue=%s", info.ID, info.Queue)

	return dto.ToWorkspaceDTO(*workspace), nil
}

func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, workspaceId string) error {
	id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		return err
	}

	workspace, err := s.workspaceRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.workspaceRepository.DeleteWorkspace(ctx, workspaceId); err != nil {
		return err
	}

	task, err := tasks.NewTerminateWorkspaceTask(id, workspace.UserID)
	if err != nil {
		return fmt.Errorf("failed to terminate workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return fmt.Errorf("failed to enqueue workspace terminating task: %w", err)
	}
	log.Printf("✅ Enqueued workspace terminating task: ID=%s queue=%s", info.ID, info.Queue)

	return nil
}

func (s *WorkspaceService) UpdateWorkspaceURL(ctx context.Context, workspaceID uint64, url string) error {
	return s.workspaceRepository.UpdateURL(ctx, workspaceID, url)
}

func (s *WorkspaceService) StartWorkspace(ctx context.Context, req requests.WorkspaceActionRequest) (bool, error) {
	err := s.UpdateWorkspaceStatus(ctx, req.ID, enums.WorkspaceStatusStarting, "Workspace starting is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	task, err := tasks.NewStartWorkspaceTask(req.ID, req.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to start workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return false, fmt.Errorf("failed to enqueue workspace starting task: %w", err)
	}
	log.Printf("✅ Enqueued workspace starting task: ID=%s queue=%s", info.ID, info.Queue)

	return true, nil
}

func (s *WorkspaceService) StopWorkspace(ctx context.Context, req requests.WorkspaceActionRequest) (bool, error) {
	err := s.UpdateWorkspaceStatus(ctx, req.ID, enums.WorkspaceStatusStopping, "Workspace stopping is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	task, err := tasks.NewStopWorkspaceTask(req.ID, req.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to stop workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return false, fmt.Errorf("failed to enqueue workspace stopping task: %w", err)
	}
	log.Printf("✅ Enqueued workspace stopping task: ID=%s queue=%s", info.ID, info.Queue)

	return true, nil
}

func (s *WorkspaceService) RestartWorkspace(ctx context.Context, req requests.WorkspaceActionRequest) (bool, error) {
	err := s.UpdateWorkspaceStatus(ctx, req.ID, enums.WorkspaceStatusRestarting, "Workspace restarting is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	task, err := tasks.NewRestartWorkspaceTask(req.ID, req.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to restart workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return false, fmt.Errorf("failed to enqueue workspace restarting task: %w", err)
	}
	log.Printf("✅ Enqueued workspace restarting task: ID=%s queue=%s", info.ID, info.Queue)

	return true, nil
}

func (s *WorkspaceService) RebuildWorkspace(ctx context.Context, req requests.WorkspaceActionRequest) (bool, error) {
	err := s.UpdateWorkspaceStatus(ctx, req.ID, enums.WorkspaceStatusRebuilding, "Workspace rebuilding is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	task, err := tasks.NewRebuildWorkspaceTask(req.ID, req.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to rebuild workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return false, fmt.Errorf("failed to enqueue workspace rebuilding task: %w", err)
	}
	log.Printf("✅ Enqueued workspace rebuilding task: ID=%s queue=%s", info.ID, info.Queue)

	return true, nil
}

func (s *WorkspaceService) TerminateWorkspace(ctx context.Context, req requests.WorkspaceActionRequest) (bool, error) {
	err := s.UpdateWorkspaceStatus(ctx, req.ID, enums.WorkspaceStatusTerminating, "Workspace terminating is waiting for processing")
	if err != nil {
		log.Printf("Failed to update workspace status: %v", err)
	}

	task, err := tasks.NewTerminateWorkspaceTask(req.ID, req.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to terminate workspace job: %w", err)
	}

	info, err := s.asynqClient.Enqueue(task, asynq.Unique(5*time.Minute))
	if err != nil {
		return false, fmt.Errorf("failed to enqueue workspace terminating task: %w", err)
	}
	log.Printf("✅ Enqueued workspace terminating task: ID=%s queue=%s", info.ID, info.Queue)

	return true, nil
}

func (s *WorkspaceService) UpdateWorkspaceStatus(ctx context.Context, workspaceID uint64, newStatus enums.WorkspaceStatus, message string) error {
	if err := s.workspaceRepository.UpdateStatus(ctx, workspaceID, string(newStatus)); err != nil {
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	event := &models.WorkspaceStatusEvent{
		WorkspaceID: workspaceID,
		Status:      newStatus,
		Message:     message,
	}
	if err := s.workspaceStatusEventRepository.Create(ctx, event); err != nil {
		return fmt.Errorf("failed to create status event: %w", err)
	}

	workspace, _ := s.GetWorkspaceIncludingDeleted(ctx, workspaceID)

	body := dto.Message{
		EventType: constants.WorkspaceStatus,
		Channel:   fmt.Sprintf("workspace_%d_status", workspaceID),
		Data: map[string]interface{}{
			"status": newStatus,
			"url":    workspace.URL,
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}
	s.publisherService.Publish(constants.CLUSTERIX_CODE_V1_EXCHANGE, constants.WORKSPACE_LOG_HANDLER_QUEUE, payload)

	return nil
}

func (s *WorkspaceService) RunWorkspaceAction(
	ctx context.Context,
	workspaceID uint64,
	userID uint64,
	onSuccess func(message string, logType string) error,
	onFailure func(message string, logType string) error,
	action constants.WorkspaceAction,
) error {
	workspace, err := s.workspaceRepository.GetByIDIncludingDeleted(ctx, workspaceID)
	if err != nil {
		wrappedErr := fmt.Errorf("workspace not found: %w", err)
		if onFailure != nil {
			_ = onFailure(err.Error(), "")
		}
		return wrappedErr
	}

	devpodWorkspaceDTO := dto.DevpodWorkspace{
		AccessToken:        workspace.GitPersonalAccessToken.Token,
		RepositoryUrl:      workspace.Repository.RepositoryURL,
		DevpodWorkspaceId:  workspace.ID,
		DevpodWorkspaceIde: "openvscode",
		AWSInstanceType:    workspace.Repository.MachineConfig.InstanceType,
		UserId:             userID,
	}

	switch action {
	case constants.ActionStart:
		aws.CreateARecord(workspace.Fingerprint)
		err = s.devpod.StartWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
	case constants.ActionStop:
		err = s.devpod.StopWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
	case constants.ActionRestart:
		err = s.devpod.RestartWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
	case constants.ActionRebuild:
		err = s.devpod.RebuildWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
	case constants.ActionTerminate:
		err = s.devpod.TerminateWorkspace(ctx, devpodWorkspaceDTO, onSuccess, onFailure)
	default:
		err = fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		wrappedErr := fmt.Errorf("devpod %s failed: %w", action, err)
		if onFailure != nil {
			_ = onFailure(err.Error(), "")
		}
		return wrappedErr
	}

	return nil
}

func (s *WorkspaceService) GenerateFingerprint(title string, userId uint64, organizationId uint32) string {
	hasher := sha256.New()

	hasher.Write([]byte(title))
	hasher.Write([]byte(fmt.Sprintf("%d", userId)))
	hasher.Write([]byte(fmt.Sprintf("%d", organizationId)))
	hasher.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	return hex.EncodeToString(hasher.Sum(nil))[:16]
}
