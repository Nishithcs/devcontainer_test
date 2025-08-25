package jobs

import (
	"clusterix-code/internal/api/requests"
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/services"
	"clusterix-code/internal/tasks"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	"strconv"
)

func NewStartWorkspaceTask(workspaceID uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(tasks.StartWorkspacePayload{
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TaskStartWorkspace, payload), nil
}

func HandleStartWorkspaceTask(ctx context.Context, t *asynq.Task, workspaceSvc *services.WorkspaceService,
	publisherSvc *services.PublisherService, workspaceConfigSvc *services.WorkspaceConfigService) error {
	var p tasks.StartWorkspacePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusStarting, "Workspace is starting by worker"); err != nil {
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	workspace, _ := workspaceSvc.GetWorkspace(ctx, p.WorkspaceID)

	successCallback := func(message string, logType string) error {

		if machineName, found := devpodParser.ExtractMachineName(message); found {
			fmt.Println("Machine Name:", machineName)
			var workspaceConfigRequest requests.UpdateWorkspaceConfigRequest
			workspaceConfigRequest.ID = workspace.WorkspaceConfig.ID
			workspaceConfigRequest.DevpodMachine = machineName
			_, err := workspaceConfigSvc.UpdateWorkspaceConfig(ctx, workspaceConfigRequest)
			if err != nil {
				log.Printf("Failed to set machine name: %v\n", err)
			}
		}


		if rawURL, found := devpodParser.ExtractURL(message); found {
			fmt.Println("Workspace URL:", rawURL)

			internalPort := devpodParser.ExtractPortFromURL(rawURL)

			mapping, err := services.MapWorkspace(strconv.FormatUint(p.WorkspaceID, 10), internalPort)
			if err != nil {
				log.Printf("Failed to map workspace port: %v\n", err)
				return nil
			}

			if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusRunning, "Workspace is running by worker"); err != nil {
				fmt.Errorf("failed to update workspace status: %w", err)
			}

			worker_base_url := os.Getenv("REVERSE_PROXY_BASE_URL")
			publicURL := fmt.Sprintf("%s.%s/?folder=/workspaces/%d", workspace.Fingerprint, worker_base_url, workspace.ID)

			if err := workspaceSvc.UpdateWorkspaceURL(ctx, p.WorkspaceID, publicURL); err != nil {
				log.Printf("Failed to update workspace URL: %v\n", err)
			}

			var workspaceConfigRequest requests.UpdateWorkspaceConfigRequest
			workspaceConfigRequest.ID = workspace.WorkspaceConfig.ID
			workspaceConfigRequest.WorkerPort = mapping.ExternalPort
			_, err = workspaceConfigSvc.UpdateWorkspaceConfig(ctx, workspaceConfigRequest)
			if err != nil {
				log.Printf("Failed to set worker port: %v\n", err)
			}
		}

		if err := PublishLogMessage(publisherSvc, p.WorkspaceID, message); err != nil {
			log.Printf("Failed to publish log message: %v", err)
		}

		return nil
	}

	failureCallback := func(message string, logType string) error {
		fmt.Println(message)

		if err := PublishLogMessage(publisherSvc, p.WorkspaceID, message); err != nil {
			log.Printf("Failed to publish log message: %v", err)
		}

		if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusFailed, "Workspace is failed by worker"); err != nil {
			fmt.Errorf("failed to update workspace status: %w", err)
		}

		return nil
	}

	if err := workspaceSvc.RunWorkspaceAction(
		ctx,
		p.WorkspaceID,
		p.UserID,
		successCallback,
		failureCallback,
		constants.ActionStart,
	); err != nil {
		return fmt.Errorf("workspace creation failed: %w", err)
	}

	return nil
}
