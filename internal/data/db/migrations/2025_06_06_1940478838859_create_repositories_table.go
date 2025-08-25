package migrations

type CreateRepositoriesTable struct {
	BaseMigration
	Name string
}

func (m *CreateRepositoriesTable) UpSql() string {
	return `CREATE TABLE repositories (
		id BIGSERIAL PRIMARY KEY,
		title VARCHAR(255),
		machine_config_id BIGINT,
		repository_url TEXT,
		status VARCHAR(20) DEFAULT 'pending',
		added_by_admin BOOLEAN,
		created_by_id BIGINT NOT NULL,
		provider_id BIGINT,
		organization_id BIGINT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP,

		FOREIGN KEY (created_by_id) REFERENCES users(id),
		FOREIGN KEY (machine_config_id) REFERENCES machine_configs(id),
		FOREIGN KEY (provider_id) REFERENCES providers(id)
	)`
}

func (m *CreateRepositoriesTable) DownSql() string {
	return "DROP TABLE IF EXISTS repositories"
}

func (m *CreateRepositoriesTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1940478838859_create_repositories_table"
}
