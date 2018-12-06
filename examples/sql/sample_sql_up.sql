CREATE TABLE  IF NOT EXISTS pika (
		id_migration SERIAL PRIMARY KEY,
		migration varchar(500) DEFAULT NULL,
		up boolean DEFAULT NULL,
		down boolean DEFAULT NULL,
		execute_up timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		execute_down timestamp DEFAULT NULL,
		statements text
	)