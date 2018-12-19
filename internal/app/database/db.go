package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/log"
)

// DB ...
type DB interface {
	Exec(sql string) error
	Save(schema Schema) error
	GetByMigrationKey(key string) error
	GetEntity(tableName string) ([]StructWithTablenName, error)
}

// StructWithTablenName ...
type StructWithTablenName struct {
	Models     string
	TableName  string
	StructName string
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

// GetEntity ...
func (s *store) GetEntity(tableName string) ([]StructWithTablenName, error) {
	query := getReverse(tableName)
	// log.Info(query)
	rows, err := s.db.Query(query)
	if err != nil {
		return []StructWithTablenName{}, err
	}
	structs := make([]StructWithTablenName, 0)
	for rows.Next() {
		var tmpl string
		var tableName string
		var structName string
		rows.Scan(&tmpl, &structName, &tableName)

		strc := StructWithTablenName{
			Models:     tmpl,
			TableName:  tableName,
			StructName: structName,
		}
		structs = append(structs, strc)
	}

	return structs, nil
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
		log.Error("failed when query get migration ", err)
		return
	}

	i := 0
	for rows.Next() {
		i++
	}

	if i > 0 {
		err = errors.New("migration already exist")
		log.Error(err)
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
		log.Error(err)
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

	log.Info("*********************************************************")
	log.Info(sql)
	log.Info("*********************************************************")

	err = rows.Close()

	return
}

func getReverse(tableName string) string {
	queryAll := ``
	if strings.ToLower(tableName) != strings.ToLower("all") {
		queryAll = fmt.Sprintf(` AND table_name = '%s' `, tableName)
	}
	return fmt.Sprintf(constant.ReversePG, "`", "`", queryAll)
}
