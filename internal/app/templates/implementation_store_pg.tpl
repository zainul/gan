package postgres

import (
	"strings"
	"{{.Package}}/internal/entity"
	"github.com/jinzhu/gorm"
)

type {{.KeyLowerCase}}Store struct {
	DB *gorm.DB
}

//New{{.Name}}Store ...
func New{{.Name}}Store(conn *gorm.DB) *{{.KeyLowerCase}}Store {
	return &{{.KeyLowerCase}}Store{
		DB: conn,
	}
}

func (s *{{.KeyLowerCase}}Store) Insert(e entity.{{.Name}}) error {
	return s.DB.Create(&e).Error
}

func (s *{{.KeyLowerCase}}Store) Update(
	field string, 
	val interface{}, 
	valuesToUpdate map[string]interface{},
) error {
	query := "UPDATE {{.KeyLowerCase}} SET "
	vals := make([]interface{}, 0)
	queryField := make([]string, 0)

	for key, valmap := range valuesToUpdate {
		vals = append(vals, valmap)

		queryField = append(queryField, key+"=?")
	}

	query = query + strings.Join(queryField, " ,") + " WHERE " + field + "= ?"
	vals = append(vals, val)

	return s.DB.Exec(query, vals...).Error
}

func (s *{{.KeyLowerCase}}Store) UpdateBy(e entity.{{.Name}}, param map[string]interface{}) error {
	return s.DB.Table(e.TableName()).Where(param).Update(&e).Error
}

func (s *{{.KeyLowerCase}}Store) Delete(id interface{}) error {
	return nil
}

func (s *{{.KeyLowerCase}}Store) RawQuery(sql string, result interface{}) error {
	return s.DB.Raw(sql).Find(result).Error
}

func (s *{{.KeyLowerCase}}Store) RawExec(query string, param ...interface{}) error {
	return s.DB.Exec(query, param...).Error
}

func (s *{{.KeyLowerCase}}Store) {{.Name}}By(field string,value interface{}) ([]entity.{{.Name}}, error) {
	result := make([]entity.{{.Name}}, 0)
	err := s.DB.Where(map[string]interface{}{field: value}).Find(&result).Error
	return result, err
}
