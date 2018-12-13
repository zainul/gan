package main

import "github.com/zainul/gan/pkg/migration"

type {{ .Key }} struct {
	migration.Migration
}

func init() {
	// will be some migration with up and down feature
	m := &{{ .Key }}{}
	migration.Register("{{ .Key }}", m)
}

func (m *{{ .Key }}) Up() {
	m.SQLFromFile(
	// `
	// create table bla bla
	// `
	)
}

func (m *{{ .Key }}) Down() {
	m.SQLFromFile(
	// `
	// drop table bla bla
	// `
	)
}
