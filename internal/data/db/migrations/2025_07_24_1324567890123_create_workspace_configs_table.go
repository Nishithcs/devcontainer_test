package migrations

type CreateWorkspaceConfigsTable struct {
	BaseMigration
	Name string
}

func (m *CreateWorkspaceConfigsTable) UpSql() string {
	return `
	CREATE TABLE workspace_configs (
		id BIGSERIAL PRIMARY KEY,
		workspace_id BIGINT,
		devpod_machine VARCHAR(100),
		worker_name VARCHAR(20),
		worker_ip VARCHAR(20),
		worker_port INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,

		FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
	)`
}

func (m *CreateWorkspaceConfigsTable) DownSql() string {
	return `DROP TABLE IF EXISTS workspace_configs`
}

func (m *CreateWorkspaceConfigsTable) GetName() string {
	return "2025_07_24_1324567890123_create_workspace_configs_table"
}
