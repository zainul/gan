package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/internal/app/constant"
)

// DB ...
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
	db     *sql.DB
	gormDB *gorm.DB
}

// NewDB ...
func NewDB(sqlconn *sql.DB) DB {
	return &store{
		db: sqlconn,
	}
}

// GetByMigrationKey ...
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

// Save ...
func (s *store) Save(schema Schema) (err error) {
	stmt, err := s.db.Prepare(constant.InsertTablePG)

	if err != nil {
		return errors.New("failed to prepare the query")
	}

	// migration, up, down, execute_up, execute_down, statements
	_, err = stmt.Exec(schema.Migration, schema.Up, schema.Down, schema.ExecuteUp, schema.ExecuteDown, schema.Statement)

	if err != nil {
		fmt.Println(err)
	}

	err = stmt.Close()

	return
}

// Exec ...
func (s *store) Exec(sql string) (err error) {

	rows, err := s.db.Query(sql)
	if err != nil {
		return err
	}

	fmt.Println("*********************************************************")
	fmt.Println(sql)
	fmt.Println("*********************************************************")

	err = rows.Close()

	return
}
