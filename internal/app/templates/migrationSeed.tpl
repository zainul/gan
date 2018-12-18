package main

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/pkg/seed"
)

// {{ .Key }} ...
{{.StructTemplate}}

// TableName ...
func (h *{{ .StructNameLower }}) TableName() string {
	return "{{.TableName}}"
}

// {{ .StructNameLower }} ...
type {{ .StructNameLower }} struct {
	db *gorm.DB
}

// Create is method creation for seed
func (s *{{ .StructNameLower }}) Create(v interface{}) error {
	{{ .StructNameLower }} := {{ .StructName }}{}

	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = json.Unmarshal(bytes, &{{ .StructNameLower }})

	if err != nil {
		return err
	}

	return s.db.Create({{ .StructNameLower }}).Error
}

// New{{ .StructName }} ...
func New{{ .StructName }}(db *gorm.DB) (seed.Store, []{{ .StructName }}) {
	arr := make([]{{ .StructName }}, 0)
	return &{{ .StructNameLower }}{db}, arr
}
