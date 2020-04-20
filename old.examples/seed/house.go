package main

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/internal/app"
)

// House ...
type House struct {
	Name string `json:"name" gorm:"column:name;"`
	No   string `json:"no" gorm:"column:no;"`
}

// TableName ...
func (h *House) TableName() string {
	return "house"
}

type house struct {
	db *gorm.DB
}

func (s *house) Create(v interface{}) error {
	house := House{}

	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = json.Unmarshal(bytes, &house)

	if err != nil {
		return err
	}

	return s.db.Create(house).Error
}

// NewHouse ...
func NewHouse(db *gorm.DB) (app.Store, []House) {
	houses := make([]House, 0)
	return &house{db}, houses
}
