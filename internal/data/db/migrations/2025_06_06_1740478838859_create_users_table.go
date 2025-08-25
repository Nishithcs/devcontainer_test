package migrations

type CreateUsersTable struct {
	BaseMigration
	Name string
}

func (m *CreateUsersTable) UpSql() string {
	return `CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			first_name VARCHAR(255), 
			last_name VARCHAR(255),
			full_name VARCHAR(255),
			email VARCHAR(255),
			avatar TEXT,
			organization_id INT,
			is_active BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		)`
}

func (m *CreateUsersTable) DownSql() string {
	return "DROP TABLE IF EXISTS users"
}

func (m *CreateUsersTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1740478838859_create_users_table"
}
