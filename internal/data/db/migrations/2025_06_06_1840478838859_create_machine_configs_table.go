package migrations

type CreateMachineConfigsTable struct {
	BaseMigration
	Name string
}

func (m *CreateMachineConfigsTable) UpSql() string {
	return `CREATE TABLE machine_configs (
		id BIGSERIAL PRIMARY KEY,
		instance_type VARCHAR(255) UNIQUE NOT NULL,
		category VARCHAR(255) NOT NULL,
		cpu_cores INTEGER NOT NULL,
		memory_gb DOUBLE PRECISION NOT NULL,
		storage_type VARCHAR(255),
		storage_size_gb INTEGER,
		network_performance VARCHAR(255),
		architecture VARCHAR(255),
		hypervisor VARCHAR(255),
		generation VARCHAR(255),
		enhanced_networking BOOLEAN,
		gpu VARCHAR(255),
		additional_features TEXT
	)`
}

func (m *CreateMachineConfigsTable) DownSql() string {
	return "DROP TABLE IF EXISTS machine_configs"
}

func (m *CreateMachineConfigsTable) GetName() string {
	// don't change this after the migration is applied
	return "2025_06_06_1740478838859_create_machine_configs_table"
}
