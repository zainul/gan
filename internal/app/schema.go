package app

import (
	"time"
)

// Schema is represent table migrations to save history of migrations db
type Schema struct {
	Migration   string     `json:"migration"`
	Up          bool       `json:"up"`
	Down        bool       `json:"down"`
	ExecuteUp   *time.Time `json:"execute_up"`
	ExecuteDown *time.Time `json:"execute_down"`
	Statement   string
}
