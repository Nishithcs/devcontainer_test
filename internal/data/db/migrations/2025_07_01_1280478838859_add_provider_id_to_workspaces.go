package migrations

type AddProviderIdToWorkspaces struct {
	BaseMigration
	Name string
}

func (m *AddProviderIdToWorkspaces) UpSql() string {
	return `
		ALTER TABLE workspaces
		ADD COLUMN provider_id BIGINT,
		ADD CONSTRAINT workspaces_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES providers(id)
	`
}

func (m *AddProviderIdToWorkspaces) DownSql() string {
	return `
		ALTER TABLE workspaces
		DROP CONSTRAINT IF EXISTS workspaces_provider_id_fkey,
		DROP COLUMN IF EXISTS provider_id
	`
}

func (m *AddProviderIdToWorkspaces) GetName() string {
	return "2025_07_01_1280478838859_add_provider_id_to_workspaces"
}
