package constant

const (
	// MigrationTablePG ...
	MigrationTablePG = `
	CREATE TABLE IF NOT EXISTS migrations (
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
	ReversePG     = `
	WITH models AS (
		WITH data AS (
		  SELECT
			  replace(initcap(table_name::text), '_', '') table_name,
			  replace(initcap(column_name::text), '_', '') column_name,
			  table_name as orig_table_name,
			  CASE data_type
			  WHEN 'timestamp without time zone' THEN 'time.Time'
			  WHEN 'timestamp with time zone' THEN 'time.Time'
			  WHEN 'boolean' THEN 'bool'
			  WHEN 'bigint' THEN 'int64'
			  WHEN 'bigserial' THEN 'int64'
			  WHEN 'double precision' THEN 'float64'
			  WHEN 'bigserial' THEN 'int64'
			  WHEN 'integer' THEN 'int'
			  WHEN 'numeric' THEN 'float64'
			  WHEN 'real' THEN 'float64'
			  WHEN 'smallint' THEN 'int64'
			  WHEN 'smallserial' THEN 'int64'
			  -- add your own type converters as needed or it will default to 'string'
			  ELSE 'string'
			  END AS type_info,
			  '%sjson:"' || column_name ||'" gorm:"column:'|| column_name ||';"%s' AS annotation
		  FROM information_schema.columns
		  WHERE table_schema IN ('public')
		  %s
		  ORDER BY table_schema, table_name, ordinal_position
		)
		  SELECT orig_table_name, table_name, STRING_AGG(E'\t' || column_name || E'\t' || type_info || E'\t' || annotation, E'\n') fields
		  FROM data
		  GROUP BY table_name, orig_table_name
	  )
	  SELECT 'type ' || table_name || E' struct {\n' || fields || E'\n}' models, table_name, orig_table_name
	
	FROM models ORDER BY 1
	`
)
