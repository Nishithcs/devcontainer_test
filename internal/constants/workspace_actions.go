package constants

type WorkspaceAction string

const (
	ActionStart     WorkspaceAction = "start"
	ActionStop      WorkspaceAction = "stop"
	ActionRestart   WorkspaceAction = "restart"
	ActionRebuild   WorkspaceAction = "rebuild"
	ActionTerminate WorkspaceAction = "terminate"
)
