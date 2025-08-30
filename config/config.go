package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"backend-task/internal/db"
)

type Config struct {
	DB_DSN    string
	DB_Driver string // "postgres"
}

var DB *gorm.DB

func ConnectToDatabase() {

	var driver = "postgres"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		GetEnv("DB_HOST", os.Getenv("DB_HOST")),
		GetEnv("DB_USER", os.Getenv("DB_USER")),
		GetEnv("DB_PASSWORD", os.Getenv("DB_PASSWORD")),
		GetEnv("DB_NAME", os.Getenv("DB_NAME")),
		GetEnv("DB_PORT", os.Getenv("DB_PORT")),
		GetEnv("DB_SSLMODE", os.Getenv("DB_SSLMODE")),
	)

	db, err := db.Open(dsn, driver)
	if err != nil {
		log.Fatal("Failed To Connect To Database :( .. ", err)
	}

	DB = db
	fmt.Println("Database Connected Successfully ^_^ .. ")
}

func GetEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
