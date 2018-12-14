package migration

import (
	"fmt"

	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/io"
)

const (
	// StatusUp ...
	StatusUp = constant.StatusUp
	// StatusDown ...
	StatusDown = constant.StatusDown
)

// Migration ..
type Migration struct {
	sql          string
	unixNanoTime float64
	key          string
}

// Register ...
func Register(name string, m app.Migrationer) {
	err := app.SetMigration(name, m)

	if err != nil {
		fmt.Println("Failed to register")
	}

	return
}

// SQL ...
func (m *Migration) SQL(sql string) {
	m.sql = sql
}

// Exec ...
func (m *Migration) Exec(status string) {
	if status == constant.StatusUp {
		app.SetExec()
	} else if status == constant.StatusDown {
		fmt.Println("status down", m.sql)
	}
}

// SQLFromFile ...
func (m *Migration) SQLFromFile(path string) {
	byteData, err := io.OpenFile(path)

	if err != nil {
		fmt.Println("Failed to open file ", err)
		return
	}

	m.sql = string(byteData)
}

// GetSQL ...
func (m *Migration) GetSQL() string {
	return m.sql
}
