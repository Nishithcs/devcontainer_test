package migrations

type AddProviderIDToWorkspacesTable struct {
	BaseMigration
	Name string
}

func (m *AddProviderIDToWorkspacesTable) UpSql() string {
	return `
		ALTER TABLE workspaces ADD COLUMN provider_id BIGINT;
		ALTER TABLE workspaces ADD CONSTRAINT workspaces_provider_id_fkey
			FOREIGN KEY (provider_id) REFERENCES providers(id);
	`
}

func (m *AddProviderIDToWorkspacesTable) DownSql() string {
	return `
		ALTER TABLE workspaces DROP CONSTRAINT IF EXISTS workspaces_provider_id_fkey;
		ALTER TABLE workspaces DROP COLUMN IF EXISTS provider_id;
	`
}

func (m *AddProviderIDToWorkspacesTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_26_1745000000000_add_provider_id_to_workspaces_table"
}
