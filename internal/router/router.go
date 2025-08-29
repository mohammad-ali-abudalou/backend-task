package router

import (
	handlers "backend-task/internal/handlers"
	services "backend-task/internal/services"

	"backend-task/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	// Wire Layers :
	repository.NewUserRepository(db)
	service := services.NewUserService(db, repository.NewUserRepository(db), repository.NewGroupRepository(db))
	handler := handlers.NewUserHandler(service)

	r := gin.Default()

	// Routes :
	r.POST("/users", handler.CreateUser)
	r.GET("/users/:id", handler.GetUserByID)
	r.PATCH("/users/:id", handler.UpdateUser)
	r.GET("/users", handler.QueryUsers) // Supports Optional Group Filter.

	return r
}
