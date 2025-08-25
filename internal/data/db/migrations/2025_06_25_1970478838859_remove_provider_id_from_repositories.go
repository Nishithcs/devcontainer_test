package migrations

type RemoveProviderIDFromRepositories struct {
	BaseMigration
	Name string
}

func (m *RemoveProviderIDFromRepositories) UpSql() string {
	return `
		ALTER TABLE repositories DROP CONSTRAINT IF EXISTS repositories_provider_id_fkey;
		ALTER TABLE repositories DROP COLUMN IF EXISTS provider_id;
	`
}

func (m *RemoveProviderIDFromRepositories) DownSql() string {
	return `
		ALTER TABLE repositories ADD COLUMN provider_id BIGINT;
		ALTER TABLE repositories ADD CONSTRAINT repositories_provider_id_fkey
			FOREIGN KEY (provider_id) REFERENCES providers(id);
	`
}

func (m *RemoveProviderIDFromRepositories) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_26_1730000000000_remove_provider_id_from_repositories"
}
