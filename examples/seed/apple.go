package main

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/internal/app"
)

// Apple ...
type Apple struct {
	//Name string `json:"name" gorm:"column:name;"`
}

// TableName ...
func (h *Apple) TableName() string {
	return "your_table_name"
}

type apple struct {
	db *gorm.DB
}

func (s *apple) Create(v interface{}) error {
	apple := Apple{}

	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = json.Unmarshal(bytes, &apple)

	if err != nil {
		return err
	}

	return s.db.Create(apple).Error
}

// NewApple ...
func NewApple(db *gorm.DB) (app.Store, []Apple) {
	arr := make([]apple, 0)
	return &house{db}, arr
}
