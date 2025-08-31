package app

import (
	"backend-task/internal/db"
	"backend-task/internal/router"
	"backend-task/internal/utils"
	"fmt"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Container struct {
	Server *router.Server
}

func InitializeContainer() *Container {

	// Initial DB :
	conn := db.InitDB()

	// Setup Routers :
	route := router.SetupRouters(conn)

	// Swagger Endpoint :
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := ":8080"
	if environmentVariableNamed := os.Getenv("HTTP_ADDR"); environmentVariableNamed != "" {
		addr = environmentVariableNamed
	}

	utils.Info(fmt.Sprintf("Listening On %s", addr))
	if err := route.Run(addr); err != nil {
		utils.Fatal(err.Error())
	}

	return &Container{Server: route}
}
