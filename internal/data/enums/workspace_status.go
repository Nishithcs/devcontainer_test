package enums

type WorkspaceStatus string

const (
	WorkspaceStatusPending     WorkspaceStatus = "pending"
	WorkspaceStatusStarting    WorkspaceStatus = "starting"
	WorkspaceStatusCreating    WorkspaceStatus = "creating"
	WorkspaceStatusRunning     WorkspaceStatus = "running"
	WorkspaceStatusStopping    WorkspaceStatus = "stopping"
	WorkspaceStatusStopped     WorkspaceStatus = "stopped"
	WorkspaceStatusRestarting  WorkspaceStatus = "restarting"
	WorkspaceStatusRebuilding  WorkspaceStatus = "rebuilding"
	WorkspaceStatusTerminating WorkspaceStatus = "terminating"
	WorkspaceStatusTerminated  WorkspaceStatus = "terminated"
	WorkspaceStatusFailed      WorkspaceStatus = "failed"
)
