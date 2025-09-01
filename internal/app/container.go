package app

import (
	"backend-task/internal/db"
	"backend-task/internal/router"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Container struct {
	Server *router.Server
}

// InitializeContainer Builds And Wires Dependencies But Does NOT Start The Server.
func InitializeContainer() *Container {

	// Initialize DB
	connection := db.InitDB()

	// Setup Routers :
	route := router.SetupRouters(connection)

	// Add Swagger Endpoint :
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Container{Server: route}
}
