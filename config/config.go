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
	DB_Driver string // "postgres" or "sqlite".
}

var DB *gorm.DB

func ConnectDatabase() {

	var driver = "postgres" // "postgres" or "sqlite".

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", os.Getenv("DB_HOST")),
		getEnv("DB_USER", os.Getenv("DB_USER")),
		getEnv("DB_PASSWORD", os.Getenv("DB_PASSWORD")),
		getEnv("DB_NAME", os.Getenv("DB_NAME")),
		getEnv("DB_PORT", os.Getenv("DB_PORT")),
	)

	db, err := db.Open(dsn, driver)
	if err != nil {
		log.Fatal("Failed To Connect To Database : ", err)
	}

	DB = db
	fmt.Println("Database Connected Successfully ^_^ .. ")
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
