package jobs

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/services"
	"clusterix-code/internal/tasks"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func NewStopWorkspaceTask(workspaceID uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(tasks.StopWorkspacePayload{
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TaskStopWorkspace, payload), nil
}

func HandleStopWorkspaceTask(ctx context.Context, t *asynq.Task, workspaceSvc *services.WorkspaceService, publisherSvc *services.PublisherService) error {
	var p tasks.StopWorkspacePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusStopping, "Workspace stopping started by worker"); err != nil {
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	successCallback := func(message string, logType string) error {
		if devpodParser.IsSuccessStop(message) {
			if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusStopped, "Workspace is stopped by worker"); err != nil {
				fmt.Errorf("failed to update workspace status: %w", err)
			}
		}

		fmt.Println(message)

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
		constants.ActionStop,
	); err != nil {
		return fmt.Errorf("workspace creation failed: %w", err)
	}

	return nil
}
