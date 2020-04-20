package seed

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/zainul/gan/internal/app"
	"github.com/zainul/gan/internal/constant"
	"github.com/zainul/gan/internal/log"
)

type Store interface {
	Create(v interface{}) error
}

func Seed(path string, store Store, value ...interface{}) {
	app.Seed(path, store, value)
}

func GetDB() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv(constant.CONNDB))
	if err != nil {
		log.Error("failed to get instance")
		return nil
	}

	return db
}
