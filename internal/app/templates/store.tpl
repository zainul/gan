package repository

import "{{.Package}}/internal/entity"

// {{.Name}} ...
type {{.Name}} interface {
	Insert(p entity.{{.Name}}) error
	Update(field string, val interface{}, valuesToUpdate map[string]interface{}) error
	UpdateBy(e entity.{{.Name}}, param map[string]interface{}) error
	Delete(id interface{}) error
	{{.Name}}By(field string,value interface{}) ([]entity.{{.Name}}, error)
	RawQuery(sql string, result interface{}) error
	RawExec(query string, param ...interface{}) error
}