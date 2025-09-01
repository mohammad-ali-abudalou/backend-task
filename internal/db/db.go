package db

import (
	"fmt"
	"os"

	"backend-task/internal/config"
	"backend-task/internal/user/models"
	"backend-task/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB Initializes The Database Connection And Runs Migrations If Enabled.
func InitDB() *gorm.DB {

	driverName := config.GetEnv(utils.DSN_DRIVER_NAME, utils.DriverPostgres) // Default: postgres
	dataSourceName := buildDSN(driverName)

	gormDB, err := open(dataSourceName, driverName)
	if err != nil {
		utils.Fatal(fmt.Sprintf("%s: %v", utils.ErrFailedConnectDatabase, err))
	}

	// Auto Migrate Schema ( Can Be Toggled Via Env ) :
	if config.GetEnv(utils.AUTO_MIGRATE, "true") == "true" {
		autoMigrate(gormDB)
		utils.Info(utils.ErrDatabaseSchemaMigratedSuccessfully.Error())
	}

	utils.Info(utils.ErrDatabaseConnected.Error())
	DB = gormDB
	return gormDB
}

func open(dsn, driverName string) (*gorm.DB, error) {

	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}

	switch driverName {
	case utils.DriverPostgres:
		return gorm.Open(postgres.Open(dsn), gcfg)

	case utils.DriverSqlite:
		return gorm.Open(sqlite.Open(dsn), gcfg)

	default:
		return nil, fmt.Errorf("%s: %s", utils.ErrUnsupportedDBDriver, driverName)

	}
}

func autoMigrate(db *gorm.DB) {

	if err := db.AutoMigrate(&models.User{}, &models.Group{}); err != nil {

		utils.Fatal(fmt.Sprintf("%s: %v", utils.ErrMigrationFailed, err))
	}
}

func buildDSN(driver string) string {

	if driver == utils.DriverSqlite {

		// Use In-Memory SQLite For Tests :
		return config.GetEnv("SQLITE_PATH", "file::memory:?cache=shared")
	}

	// Default: Postgres DSN
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.GetEnv(utils.DSN_DB_HOST, os.Getenv(utils.DSN_DB_HOST)),
		config.GetEnv(utils.DSN_DB_USER, os.Getenv(utils.DSN_DB_USER)),
		config.GetEnv(utils.DSN_DB_PASSWORD, os.Getenv(utils.DSN_DB_PASSWORD)),
		config.GetEnv(utils.DSN_DB_NAME, os.Getenv(utils.DSN_DB_NAME)),
		config.GetEnv(utils.DSN_DB_PORT, os.Getenv(utils.DSN_DB_PORT)),
		config.GetEnv(utils.DSN_DB_SSLMODE, os.Getenv(utils.DSN_DB_SSLMODE)),
	)
}
