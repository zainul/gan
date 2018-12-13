package main

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/pkg/seed"
)

// {{ .Key }} ...
type {{ .Key }} struct {
	//Name string `json:"name" gorm:"column:name;"`
}

// TableName ...
func (h *{{ .Key }}) TableName() string {
	return "your_table_name"
}

type {{ .KeyLowerCase }} struct {
	db *gorm.DB
}

func (s *{{ .KeyLowerCase }}) Create(v interface{}) error {
	{{ .KeyLowerCase }} := {{ .Key }}{}

	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = json.Unmarshal(bytes, &{{ .KeyLowerCase }})

	if err != nil {
		return err
	}

	return s.db.Create({{ .KeyLowerCase }}).Error
}

// New{{ .Key }} ...
func New{{ .Key }}(db *gorm.DB) (seed.Store, []{{ .Key }}) {
	arr := make([]{{ .Key }}, 0)
	return &{{ .KeyLowerCase }}{db}, arr
}
