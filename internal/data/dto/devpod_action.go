package dto

type DevpodWorkspace struct {
	AccessToken        string `json:"access_token"`
	RepositoryUrl      string `json:"repository_url"`
	DevpodWorkspaceId  uint64 `json:"devpod_workspace_id"`
	DevpodWorkspaceIde string `json:"devpod_workspace_ide"`
	AWSInstanceType    string `json:"aws_instance_type"`
	UserId             uint64 `json:"user_id"`
}
