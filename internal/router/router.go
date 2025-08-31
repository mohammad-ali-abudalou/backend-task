package router

import (
	handlers "backend-task/internal/handlers"
	service "backend-task/internal/services"
	services "backend-task/internal/services"

	"backend-task/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouters(db *gorm.DB) *gin.Engine {

	// Layers of Wire :
	repository.NewUserRepository(db)
	service := services.NewUserService(db, repository.NewUserRepository(db), repository.NewGroupRepository(db))
	handler := handlers.NewUserHandler(service)

	route := gin.Default()

	// Routes :
	route.POST("/users", handler.CreateUser)
	route.GET("/users/:id", handler.GetUserByID)
	route.PATCH("/users/:id", handler.UpdateUser)
	route.GET("/users", handler.QueryUsers) // Group Filter.

	return route
}

// Injecting A Mock Service :
func SetupRoutersWithService(userService service.UserService) *gin.Engine {

	route := gin.Default()

	handler := handlers.NewUserHandler(userService)

	// Routes :
	route.POST("/users", handler.CreateUser)
	route.GET("/users/:id", handler.GetUserByID)
	route.PATCH("/users/:id", handler.UpdateUser)
	route.GET("/users", handler.QueryUsers) // Group Filter.

	return route
}
