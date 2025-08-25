package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const TaskRestartWorkspace = "workspace:restart"

type RestartWorkspacePayload struct {
	WorkspaceID uint64
	UserID      uint64
}

func NewRestartWorkspaceTask(workspaceID uint64, userId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(RestartWorkspacePayload{
		WorkspaceID: workspaceID,
		UserID:      userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskRestartWorkspace, payload), nil
}
