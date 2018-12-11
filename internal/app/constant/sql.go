package constant

const (
	// MigrationTablePG ...
	MigrationTablePG = `
	CREATE TABLE  IF NOT EXISTS migrations (
		id_migration SERIAL PRIMARY KEY,
		migration varchar(500) DEFAULT NULL,
		up boolean DEFAULT NULL,
		down boolean DEFAULT NULL,
		execute_up timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		execute_down timestamp DEFAULT NULL,
		statements text
	)
	`
	// InsertTablePG ...
	InsertTablePG = `INSERT INTO migrations (migration, up, down, execute_up, execute_down, statements) values ($1, $2, $3, $4, $5, $6)`
)
