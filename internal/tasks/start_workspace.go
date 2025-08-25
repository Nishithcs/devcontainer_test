package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const TaskStartWorkspace = "workspace:start"

type StartWorkspacePayload struct {
	WorkspaceID uint64
	UserID      uint64
}

func NewStartWorkspaceTask(workspaceID uint64, userId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(StartWorkspacePayload{
		WorkspaceID: workspaceID,
		UserID:      userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskStartWorkspace, payload), nil
}
