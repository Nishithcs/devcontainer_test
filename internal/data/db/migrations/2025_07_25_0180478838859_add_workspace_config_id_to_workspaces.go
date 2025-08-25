package migrations

type AddWorkspaceConfigIdToWorkspaces struct {
	BaseMigration
	Name string
}

func (m *AddWorkspaceConfigIdToWorkspaces) UpSql() string {
	return `
		ALTER TABLE workspaces
		ADD COLUMN workspace_config_id BIGINT,
		ADD CONSTRAINT workspaces_workspace_config_id_fkey FOREIGN KEY (workspace_config_id) REFERENCES workspace_configs(id)
	`
}

func (m *AddWorkspaceConfigIdToWorkspaces) DownSql() string {
	return `
		ALTER TABLE workspaces
		DROP CONSTRAINT IF EXISTS workspaces_workspace_config_id_fkey,
		DROP COLUMN IF EXISTS workspace_config_id
	`
}

func (m *AddWorkspaceConfigIdToWorkspaces) GetName() string {
	return "2025_07_01_1280478838859_add_workspace_config_id_to_workspaces"
}
