package db

import (
	"fmt"

	"backend-task/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(dsn string, driver string) (*gorm.DB, error) {

	var (
		db  *gorm.DB
		err error
	)

	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
	switch driver {

	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), gcfg)

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), gcfg)

	default:
		return nil, fmt.Errorf("Unsupported DB Driver: %s", driver)
	}

	return db, err
}

func Migrate(db *gorm.DB) error {

	return db.AutoMigrate(&models.User{}, &models.Group{})
}
