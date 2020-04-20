package main

import "github.com/zainul/gan/internal/app"

type house_table_20181207_143748_1544168268934406236 struct {
	app.Migration
}

func init() {
	// will be some migration with up and down feature
	m := &house_table_20181207_143748_1544168268934406236{}
	app.Register("house_table_20181207_143748_1544168268934406236", m)
}

func (m *house_table_20181207_143748_1544168268934406236) Up() {
	m.SQL(
		`CREATE TABLE  IF NOT EXISTS house (
			name text,
			no text
		)`,
	)
}

func (m *house_table_20181207_143748_1544168268934406236) Down() {
	m.SQL(`DROP TABLE house`)
}
