package main

import (
	"encoding/json"
	"errors"

	{{.Dependencies}}

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/pkg/seed"
)

// {{ .Key }} ...
{{.StructTemplate}}

// TableName ...
func (h *{{ .StructName }}) TableName() string {
	return "{{.TableName}}"
}

// store{{ .StructName }} ...
type store{{ .StructName }} struct {
	db *gorm.DB
}

// Create is method creation for seed
func (s *store{{ .StructName }}) Create(v interface{}) error {
	obj{{ .StructName }} := {{ .StructName }}{}

	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = json.Unmarshal(bytes, &obj{{ .StructName }})

	if err != nil {
		return err
	}

	return s.db.Create(obj{{ .StructName }}).Error
}

// New{{ .StructName }} ...
func New{{ .StructName }}(db *gorm.DB) (seed.Store, []{{ .StructName }}) {
	arr := make([]{{ .StructName }}, 0)
	return &store{{ .StructName }}{db}, arr
}
