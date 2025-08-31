package db

import (
	"fmt"
	"log"
	"os"

	"backend-task/internal/user/models"
	"backend-task/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {

	var driverName = os.Getenv("DRIVER_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", os.Getenv("DB_HOST")),
		getEnv("DB_USER", os.Getenv("DB_USER")),
		getEnv("DB_PASSWORD", os.Getenv("DB_PASSWORD")),
		getEnv("DB_NAME", os.Getenv("DB_NAME")),
		getEnv("DB_PORT", os.Getenv("DB_PORT")),
		getEnv("DB_SSLMODE", os.Getenv("DB_SSLMODE")),
	)

	gormDB, err := open(dsn, driverName)
	if err != nil {
		log.Fatal(utils.ErrFailedConnectDatabase, err)
	}

	// Auto Migrate Schema :
	AutoMigrate(gormDB)

	fmt.Println(utils.ErrDatabaseConnectedSuccessfully)

	DB = gormDB

	return gormDB
}

func getEnv(key, fallback string) string {

	if environmentVariableNamed, exists := os.LookupEnv(key); exists {
		return environmentVariableNamed
	}

	return fallback
}

func open(dsn string, driverName string) (*gorm.DB, error) {

	var (
		db  *gorm.DB
		err error
	)

	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Error)}

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

	if err := db.AutoMigrate(&models.User{}, &models.Group{}); err != nil {

		log.Fatalf("%s : %s", utils.ErrMigrationFailed, err)
	}
}
