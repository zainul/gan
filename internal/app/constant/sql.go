package constant

const (
	MigrationTablePG = `
	CREATE TABLE  IF NOT EXISTS migrations (
		id_migration SERIAL PRIMARY KEY,
		migration varchar(500) DEFAULT NULL,
		up boolean DEFAULT NULL,
		down boolean DEFAULT NULL,
		execute_up timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		execute_down timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		statements text
	)
	`
)
