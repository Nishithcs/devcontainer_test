package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const TaskStopWorkspace = "workspace:stop"

type StopWorkspacePayload struct {
	WorkspaceID uint64
	UserID      uint64
}

func NewStopWorkspaceTask(workspaceID uint64, userId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(StopWorkspacePayload{
		WorkspaceID: workspaceID,
		UserID:      userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskStopWorkspace, payload), nil
}
