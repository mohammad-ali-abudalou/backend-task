package router

import (
	handlers "backend-task/internal/user/handlers"
	services "backend-task/internal/user/services"

	"backend-task/internal/user/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	*gin.Engine
}

func SetupRouters(db *gorm.DB) *Server {

	route := gin.Default()

	// Layers of Wire :
	repository.NewUserRepository(db)
	service := services.NewUserService(db, repository.NewUserRepository(db), repository.NewGroupRepository(db))
	handler := handlers.NewUserHandler(service)

	// Routes :
	api := route.Group("/api")
	{
		api.POST("/users", handler.CreateUser)
		api.GET("/users/:id", handler.GetUserByID)
		api.PATCH("/users/:id", handler.UpdateUser)
		api.GET("/users", handler.QueryUsers) // Group Filter.
	}

	return &Server{route}
}

// Injecting A Mock Service :
func SetupRoutersWithService(userService services.UserService) *gin.Engine {

	route := gin.Default()

	handler := handlers.NewUserHandler(userService)

	// Routes :
	route.POST("/users", handler.CreateUser)
	route.GET("/users/:id", handler.GetUserByID)
	route.PATCH("/users/:id", handler.UpdateUser)
	route.GET("/users", handler.QueryUsers) // Group Filter.

	return route
}
