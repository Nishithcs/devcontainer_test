package migrations

type CreateWorkspaceStatusEventsTable struct {
	BaseMigration
	Name string
}

func (m *CreateWorkspaceStatusEventsTable) UpSql() string {
	return `CREATE TABLE workspace_status_events (
		id BIGSERIAL PRIMARY KEY,
		workspace_id BIGINT NOT NULL,
		status VARCHAR(50) NOT NULL,
		message TEXT,
		created_at TIMESTAMP DEFAULT NOW(),

		FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
	);`
}

func (m *CreateWorkspaceStatusEventsTable) DownSql() string {
	return "DROP TABLE IF EXISTS workspace_status_events"
}

func (m *CreateWorkspaceStatusEventsTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_30_1400000000000_create_workspace_status_events_table"
}
