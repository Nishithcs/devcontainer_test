package migrations

type CreateGitPersonalAccessTokensTable struct {
	BaseMigration
	Name string
}

func (m *CreateGitPersonalAccessTokensTable) UpSql() string {
	return `CREATE TABLE git_personal_access_tokens (
		id BIGSERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		token VARCHAR(255) NOT NULL,
		user_id BIGINT NOT NULL,
		is_default BOOLEAN,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,

		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
}

func (m *CreateGitPersonalAccessTokensTable) DownSql() string {
	return "DROP TABLE IF EXISTS git_personal_access_tokens"
}

func (m *CreateGitPersonalAccessTokensTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1950478838859_create_git_personal_access_tokens_table"
}
