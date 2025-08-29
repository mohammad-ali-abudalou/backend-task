package main

import (
	"log"

	"backend-task/config"
	"backend-task/internal/models"
	"backend-task/internal/router"
)

func main() {

	// Connect DB :
	config.ConnectDatabase()

	// Auto Migrate Schema :
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration Failed : %v", err)
	}

	// Routes :
	router.SetupRouter(config.DB)
}
