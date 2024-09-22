package main

import (
	_ "MAXPUMP1/cmd/api/docs"
	"MAXPUMP1/pkg/di"
	"log"
	"os"

	routes "MAXPUMP1/pkg/api/route"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	//Loading .env file
	if os.Getenv("GO_ENV") != "test" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("error in loading .env file: %v", err)
		}
	}

	// Initialize UserHandler
	userHandler := di.InitializeUserApi()

	// Initialize AdminHandler
	adminHandler := di.InitializeAdminApi()
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.LoadHTMLGlob("pkg/assets/*")

	// Register routes for UserHandler
	router = routes.UserRoutes(router, userHandler)

	// Register routes for AdminHandler
	router = routes.AdminRoutes(router, adminHandler)
	router.Run(":8080")

}
