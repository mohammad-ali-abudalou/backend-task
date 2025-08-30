package main

import (
	"log"
	"os"

	"backend-task/config"
	"backend-task/internal/models"
	"backend-task/internal/router"

	"github.com/joho/godotenv"

	// swagger docs
	_ "backend-task/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env File Found !")
	}

	// Connect To DB :
	config.ConnectToDatabase()

	// Auto Migrate Schema :
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration Failed : %v", err)
	}

	// Routes :
	r := router.SetupRouters(config.DB)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := ":8080"
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		addr = v
	}

	log.Printf("Listening On %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
