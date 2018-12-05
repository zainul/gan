package app

import "errors"

type Migrationer interface {
	Up()
	Down()
}

var (
	migrationMap map[string]Migrationer
)

func init() {
	migrationMap = make(map[string]Migrationer)
}

type Migration struct {
	SQL string
}

func Register(name string, m Migrationer) error {
	if _, ok := migrationMap[name]; ok {
		return errors.New("already exist name:" + name)
	}
	migrationMap[name] = m
	return nil
}
