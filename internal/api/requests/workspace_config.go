package requests

type CreateWorkspaceConfigRequest struct {
	WorkspaceID   uint64 `json:"workspace_id" binding:"omitempty"`
	DevpodMachine string `json:"devpod_machine" binding:"omitempty"`
	WorkerName    string `json:"worker_name" binding:"omitempty"`
	WorkerIP      string `json:"worker_ip" binding:"omitempty"`
	WorkerPort    int    `json:"worker_port" binding:"omitempty"`
}

type UpdateWorkspaceConfigRequest struct {
	ID            uint64 `json:"id" binding:"required"`
	DevpodMachine string `json:"devpod_machine" binding:"omitempty"`
	WorkerName    string `json:"worker_name" binding:"omitempty"`
	WorkerIP      string `json:"worker_ip" binding:"omitempty"`
	WorkerPort    int    `json:"worker_port" binding:"omitempty"`
}
