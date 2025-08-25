package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const TaskTerminateWorkspace = "workspace:terminate"

type TerminateWorkspacePayload struct {
	WorkspaceID uint64
	UserID      uint64
}

func NewTerminateWorkspaceTask(workspaceID uint64, userId uint64) (*asynq.Task, error) {
	payload, err := json.Marshal(TerminateWorkspacePayload{
		WorkspaceID: workspaceID,
		UserID:      userId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTerminateWorkspace, payload), nil
}
