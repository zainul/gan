package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type DB interface {
	Exec(sql string) error
}

type store struct {
	db *sql.DB
}

func NewDB(sqlconn *sql.DB) DB {
	return &store{
		db: sqlconn,
	}
}

func (s *store) Exec(sql string) (err error) {
	stmt, err := s.db.Prepare(sql)

	if err != nil {
		fmt.Println("error when prepare query ", err)
		return errors.New("failed to prepare the query")
	}
	_, errStmt := stmt.Exec()

	defer func() {
		err = stmt.Close()
	}()

	return errStmt
}
