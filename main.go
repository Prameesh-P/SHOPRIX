package main

import (
	
	_"github.com/Prameesh-P/SHOPRIX/docs"
	"github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/Prameesh-P/SHOPRIX/routes"
	"github.com/gin-gonic/gin"
	"os"
	swaggerFiles"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initalizers.LoadEnvVariables()
	database.ConnectToDb()
	database.SyncDb()
}
// @title Gin Swagger Example API
// @version 1.0
// @description This is a Complete Ecormmerce server.
// @termsOfService http://swagger.io/terms/
// @host localhost:9000
// @BasePath /
// @schemes http

func main() {
	controllers.ImageResizing()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.New()
	//url := ginSwagger.URL("http://localhost:9000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(gin.Recovery(), middlewares.Logger())
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run()

}
