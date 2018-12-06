package main

import "github.com/zainul/gan/internal/app"

type {{ .Key }} struct {
	app.Migration
}

func init() {
	// will be some migration with up and down feature
	m := &{{ .Key }}{}
	app.Register("{{ .Key }}", m)
}

func (m *{{ .Key }}) Up() {
	// start from directory that have you setting in config.json
	m.SQLFromFile(
	// `
	// create table bla bla
	// `
	)
}

func (m *{{ .Key }}) Down() {
	// start from directory that have you setting in config.json
	m.SQLFromFile(
	// `
	// drop table bla bla
	// `
	)
}
