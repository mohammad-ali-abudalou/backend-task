package main

import (
	"log"
	"os"

	"backend-task/config"
	"backend-task/internal/models"
	"backend-task/internal/router"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env File.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env File Found !")
	}

	// Connect DB :
	config.ConnectDatabase()

	// Auto Migrate Schema :
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration Failed : %v", err)
	}

	// Routes :
	r := router.SetupRouter(config.DB)

	addr := ":8080"
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		addr = v
	}

	log.Printf("Listening On %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
