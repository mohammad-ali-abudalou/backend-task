package main

import (

    "backend-task/config"
	"backend-task/internal/models"
    "backend-task/internal/router"
	
    "fmt"
	"log"
)

func main() {
	
    fmt.Println("Project is running 🚀")
	
	// Connect DB
	config.ConnectDatabase()

	// Auto migrate schema
	config.DB.AutoMigrate(&models.User{})
	
	
	db, err := config.ConnectDB()
    if err != nil {
        log.Fatal("DB connection error:", err)
    }

    r := router.SetupRouter(db)
    log.Println("Server running on :8080")
    r.Run(":8080")
}