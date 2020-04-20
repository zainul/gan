package main

import (
	"github.com/zainul/gan/pkg/migration"
)

func main() {
	m := migration.Migration{}
	m.Exec(migration.StatusUp)
}
