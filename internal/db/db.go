package db

import (
	"fmt"
	"log"

	"backend-task/internal/models"
	"backend-task/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(dsn string, driverName string) (*gorm.DB, error) {

	var (
		db  *gorm.DB
		err error
	)

	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	switch driverName {

	case utils.DriverPostgres:
		db, err = gorm.Open(postgres.Open(dsn), gcfg)

	case utils.DriverSqlite:
		db, err = gorm.Open(sqlite.Open(dsn), gcfg)

	default:
		return nil, fmt.Errorf("%s : %s", utils.ErrUnsupportedDBDriver, driverName)
	}

	return db, err
}

func AutoMigrate(db *gorm.DB) {

	if err := db.AutoMigrate(&models.User{}); err != nil {

		log.Fatalf("%s : %s", utils.ErrMigrationFailed, err)
	}
}
