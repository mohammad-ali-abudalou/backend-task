package main

import (
	"log"

	"backend-task/internal/app"
	"backend-task/internal/config"

	_ "backend-task/docs"
)

// @title Backend Task API
// @version 1.0
// @description REST API in Go (Gin + GORM) with automatic group assignment.
// @contact.name Mohammad Ali Abu-Dalou
// @contact.email mohammadaliabudalou@example.com
// @host localhost:8080
// @BasePath /
func main() {

	// Load Values From .env File.
	config.LoadEnv()

	container := app.InitializeContainer()

	server := container.Server

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
