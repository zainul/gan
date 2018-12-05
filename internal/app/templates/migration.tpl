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
	m.SQL(
	// `
	// create table bla bla
	// `
	)
}

func (m *{{ .Key }}) Down() {
	m.SQL(
	// `
	// drop table bla bla
	// `
	)
}
