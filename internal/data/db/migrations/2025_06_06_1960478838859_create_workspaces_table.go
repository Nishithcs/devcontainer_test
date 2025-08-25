package migrations

type CreateWorkspacesTable struct {
	BaseMigration
	Name string
}

func (m *CreateWorkspacesTable) UpSql() string {
	return `CREATE TABLE workspaces (
		id BIGSERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		color VARCHAR(50),
		ide VARCHAR(50),
		repository_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		url TEXT,
		organization_id BIGINT NOT NULL,
		git_personal_access_token_id BIGINT NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		tags TEXT[],
		last_run_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,

		FOREIGN KEY (repository_id) REFERENCES repositories(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (git_personal_access_token_id) REFERENCES git_personal_access_tokens(id)
	)`
}

func (m *CreateWorkspacesTable) DownSql() string {
	return "DROP TABLE IF EXISTS workspaces"
}

func (m *CreateWorkspacesTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1960478838859_create_workspaces_table"
}
