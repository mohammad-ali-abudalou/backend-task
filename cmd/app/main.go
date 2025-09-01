package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-task/internal/app"
	"backend-task/internal/config"
	"backend-task/internal/utils"

	_ "backend-task/docs" // <- Swagger docs import
)

// @title Backend Task API
// @version 1.0
// @description REST API in Go (Gin + GORM) with automatic group assignment.
// @contact.name Mohammad Ali Abu-Dalou
// @contact.email mohammad_abudalou@hotmail.com
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// Load Values From .env File.
	config.LoadEnv()

	// Initialize Container :
	container := app.InitializeContainer()
	server := container.Server

	// Configurable Port :
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Setup Custom HTTP Server :
	customserver := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	// Start Server In Goroutine :
	go func() {

		log.Printf("Server Running On http://localhost%s", addr)
		if err := customserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server Failed To Start")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Fatal("Shutting Down Server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := customserver.Shutdown(ctx); err != nil {
		log.Fatal("Forced Shutdown Due To Timeout")
	}

	utils.Info("Server Exited Gracefully")
}
