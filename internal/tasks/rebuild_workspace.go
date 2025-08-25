package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const TaskRebuildWorkspace = "workspace:rebuild"

type RebuildWorkspacePayload struct {
	WorkspaceID uint64
	UserID      uint64
}

func NewRebuildWorkspaceTask(workspaceID uint64, userId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(RebuildWorkspacePayload{
		WorkspaceID: workspaceID,
		UserID:      userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskRebuildWorkspace, payload), nil
}
