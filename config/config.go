package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"backend-task/internal/db"
	"backend-task/pkg/utils"
)

type Config struct {
	DB_DSN    string
	DB_Driver string
}

var DB *gorm.DB

func ConnectToDatabase() {

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

	db, err := db.Open(dsn, driverName)
	if err != nil {
		log.Fatal(utils.ErrFailedConnectDatabase, err)
	}

	DB = db
	fmt.Println(utils.ErrDatabaseConnectedSuccessfully)
}

func getEnv(key, fallback string) string {

	if environmentVariableNamed, exists := os.LookupEnv(key); exists {
		return environmentVariableNamed
	}

	return fallback
}
