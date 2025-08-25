package jobs

import (
	"clusterix-code/internal/constants"
	"clusterix-code/internal/data/enums"
	"clusterix-code/internal/services"
	"clusterix-code/internal/tasks"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

func NewRebuildWorkspaceTask(workspaceID uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(tasks.RebuildWorkspacePayload{
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TaskRebuildWorkspace, payload), nil
}

func HandleRebuildWorkspaceTask(ctx context.Context, t *asynq.Task, workspaceSvc *services.WorkspaceService, publisherSvc *services.PublisherService) error {
	var p tasks.RebuildWorkspacePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if err := workspaceSvc.UpdateWorkspaceStatus(ctx, p.WorkspaceID, enums.WorkspaceStatusCreating, "Workspace creation started by worker"); err != nil {
		return fmt.Errorf("failed to update workspace status: %w", err)
	}

	successCallback := func(message string, logType string) error {
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

		return nil
	}

	if err := workspaceSvc.RunWorkspaceAction(
		ctx,
		p.WorkspaceID,
		p.UserID,
		successCallback,
		failureCallback,
		constants.ActionRebuild,
	); err != nil {
		return fmt.Errorf("workspace creation failed: %w", err)
	}

	return nil
}
