package main

import "github.com/zainul/gan/pkg/migration"

// {{ .Key }} ...
type {{ .Key }} struct {
	migration.Migration
}

func init() {
	// will be some migration with up and down feature
	m := &{{ .Key }}{}
	migration.Register("{{ .Key }}", m)
}

// Up is migration up
func (m *{{ .Key }}) Up() {
	m.SQL(
	// `
	// create table bla bla
	// `
	)
}

// Down migration down
func (m *{{ .Key }}) Down() {
	m.SQL(
	// `
	// drop table bla bla
	// `
	)
}
