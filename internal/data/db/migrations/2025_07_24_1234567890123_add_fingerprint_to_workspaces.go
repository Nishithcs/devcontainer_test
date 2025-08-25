package migrations

type AddFingerprintToWorkspaces struct {
	BaseMigration
	Name string
}

func (m *AddFingerprintToWorkspaces) UpSql() string {
	return `
		ALTER TABLE workspaces
		ADD COLUMN fingerprint VARCHAR(100) UNIQUE;
	`
}

func (m *AddFingerprintToWorkspaces) DownSql() string {
	return `
		ALTER TABLE workspaces
		DROP COLUMN IF EXISTS fingerprint;
	`
}

func (m *AddFingerprintToWorkspaces) GetName() string {
	// Don't change this after applying the migration
	return "2025_07_24_1234567890123_add_fingerprint_to_workspaces"
}
