package migrations

type CreateProvidersTable struct {
	BaseMigration
	Name string
}

func (m *CreateProvidersTable) UpSql() string {
	return `CREATE TABLE providers (
		id BIGSERIAL PRIMARY KEY,
		title VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
	)`
}

func (m *CreateProvidersTable) DownSql() string {
	return "DROP TABLE IF EXISTS providers"
}

func (m *CreateProvidersTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1750478838859_create_providers_table"
}
