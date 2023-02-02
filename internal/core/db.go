package core

import (
	"fmt"

	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(appConfig models.Config) *gorm.DB {
	var err error

	dsn := fmt.Sprintf("%s.db", appConfig.DbName)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
