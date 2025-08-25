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

func NewTerminateWorkspaceTask(workspaceID uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(tasks.TerminateWorkspacePayload{
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TaskTerminateWorkspace, payload), nil
}

func HandleTerminateWorkspaceTask(ctx context.Context, t *asynq.Task, workspaceSvc *services.WorkspaceService, publisherSvc *services.PublisherService) error {
	var p tasks.TerminateWorkspacePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusTerminating, "Workspace terminating started by worker"); err != nil {
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	successCallback := func(message string, logType string) error {
		if devpodParser.IsSuccessDelete(message) {
			if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusTerminated, "Workspace is terminated by worker"); err != nil {
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
		constants.ActionTerminate,
	); err != nil {
		return fmt.Errorf("workspace terminating failed: %w", err)
	}
	
	return nil
}
