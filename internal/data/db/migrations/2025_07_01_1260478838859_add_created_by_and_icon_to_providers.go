package migrations

type AddCreatedByAndIconToProviders struct {
	BaseMigration
}

func (m *AddCreatedByAndIconToProviders) UpSql() string {
	return `
		ALTER TABLE providers
		ADD COLUMN created_by_id BIGINT,
		ADD COLUMN icon TEXT,
		ADD CONSTRAINT fk_providers_created_by FOREIGN KEY (created_by_id) REFERENCES users(id)
	`
}

func (m *AddCreatedByAndIconToProviders) DownSql() string {
	return `
		ALTER TABLE providers
		DROP CONSTRAINT IF EXISTS fk_providers_created_by,
		DROP COLUMN IF EXISTS created_by_id,
		DROP COLUMN IF EXISTS icon
	`
}

func (m *AddCreatedByAndIconToProviders) GetName() string {
	// unique and consistent name
	return "2025_07_01_1260478838859_add_created_by_and_icon_to_providers"
}
