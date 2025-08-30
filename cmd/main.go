package main

import (
	"log"
	"os"

	"backend-task/config"
	"backend-task/internal/router"
	"backend-task/pkg/utils"

	"github.com/joho/godotenv"

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
		log.Println(utils.ErrNoEnvFileFound)
	}

	// Connect To DB :
	config.ConnectToDatabase()

	// Auto Migrate Schema :
	config.DB.AutoMigrate(config.DB)

	// Routes :
	route := router.SetupRouters(config.DB)

	// Swagger Endpoint :
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := ":8080"
	if environmentVariableNamed := os.Getenv("HTTP_ADDR"); environmentVariableNamed != "" {
		addr = environmentVariableNamed
	}

	log.Printf("Listening On %s", addr)
	if err := route.Run(addr); err != nil {
		log.Fatal(err)
	}
}
