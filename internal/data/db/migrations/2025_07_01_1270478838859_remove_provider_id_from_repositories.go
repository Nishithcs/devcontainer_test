package migrations

type RemoveProviderIdFromRepositories struct {
	BaseMigration
	Name string
}

func (m *RemoveProviderIdFromRepositories) UpSql() string {
	return `
		ALTER TABLE repositories
		DROP CONSTRAINT IF EXISTS repositories_provider_id_fkey,
		DROP COLUMN IF EXISTS provider_id
	`
}

func (m *RemoveProviderIdFromRepositories) DownSql() string {
	return `
		ALTER TABLE repositories
		ADD COLUMN provider_id BIGINT,
		ADD CONSTRAINT repositories_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES providers(id)
	`
}

func (m *RemoveProviderIdFromRepositories) GetName() string {
	return "2025_07_01_1270478838859_remove_provider_id_from_repositories"
}
