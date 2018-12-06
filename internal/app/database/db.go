package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zainul/gan/internal/app/constant"
)

type DB interface {
	Exec(sql string) error
	Save(schema Schema) error
	GetByMigrationKey(key string) error
}

// Schema is represent table migrations to save history of migrations db
type Schema struct {
	Migration   string     `json:"migration"`
	Up          bool       `json:"up"`
	Down        bool       `json:"down"`
	ExecuteUp   *time.Time `json:"execute_up"`
	ExecuteDown *time.Time `json:"execute_down"`
	Statement   string
}

type store struct {
	db *sql.DB
}

func NewDB(sqlconn *sql.DB) DB {
	return &store{
		db: sqlconn,
	}
}

func (s *store) GetByMigrationKey(key string) (err error) {
	query := `SELECT id_migration from migrations where migration = $1`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to prepare the query %v", err.Error())
	}
	rows, err := s.db.Query(query, key)

	err = stmt.Close()

	if err != nil {
		fmt.Println("failed when query get migration ", err)
		return
	}

	i := 0
	for rows.Next() {
		i++
	}

	if i > 0 {
		err = errors.New("migration already exist")
		fmt.Println(err)
		return
	}

	return
}

func (s *store) Save(schema Schema) (err error) {
	stmt, err := s.db.Prepare(constant.InsertTablePG)

	if err != nil {
		return errors.New("failed to prepare the query")
	}

	// migration, up, down, execute_up, execute_down, statements
	_, err = stmt.Exec(schema.Migration, schema.Up, schema.Down, schema.ExecuteUp, schema.ExecuteDown, schema.Statement)

	fmt.Println(err)

	err = stmt.Close()

	return
}

func (s *store) Exec(sql string) (err error) {
	stmt, err := s.db.Prepare(sql)

	if err != nil {
		fmt.Println("error when prepare query ", err)
		return errors.New("failed to prepare the query")
	}
	_, err = stmt.Exec()

	if err == nil {
		fmt.Println("*********************************************************")
		fmt.Println(sql)
		fmt.Println("*********************************************************")
	}

	err = stmt.Close()

	return
}
