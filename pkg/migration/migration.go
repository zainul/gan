package migration

import (
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/app/constant"
	"github.com/zainul/gan/internal/app/io"
	"github.com/zainul/gan/internal/app/log"
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
		log.Error("Failed to register")
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
		log.Info("status down", m.sql)
	}
}

// SQLFromFile ...
func (m *Migration) SQLFromFile(path string) {
	byteData, err := io.OpenFile(path)

	if err != nil {
		log.Error("Failed to open file ", err)
		return
	}

	m.sql = string(byteData)
}

// GetSQL ...
func (m *Migration) GetSQL() string {
	return m.sql
}
