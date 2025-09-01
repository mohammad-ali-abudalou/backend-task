package router

import (
	"backend-task/internal/user/handlers"
	"backend-task/internal/user/repository"
	services "backend-task/internal/user/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	*gin.Engine
}

// Setup Routers Builds All Routes And Dependencies :
func SetupRouters(db *gorm.DB) *Server {

	router := gin.New()

	// Middlewares :
	router.Use(gin.Logger())   // Request Logging.
	router.Use(gin.Recovery()) // Rcover From Panics.
	router.Use(cors.Default()) // CORS Enabled By Default.

	// Wire layers :
	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	userService := services.NewUserService(db, userRepo, groupRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Versioned API Routes :
	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUserByID)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.GET("/users", userHandler.QueryUsers) // Supports Group Filter.
	}

	// Health Check ( Useful For Kubernetes, Docker, etc. )
	router.GET("/health", func(context *gin.Context) {

		context.JSON(200, gin.H{"status": "ok"})
	})

	return &Server{router}
}

// For Testing With Mocks :
func SetupRoutersWithService(userService services.UserService) *gin.Engine {

	router := gin.Default()
	handler := handlers.NewUserHandler(userService)

	router.POST("/users", handler.CreateUser)
	router.GET("/users/:id", handler.GetUserByID)
	router.PUT("/users/:id", handler.UpdateUser)
	router.GET("/users", handler.QueryUsers)

	return router
}
