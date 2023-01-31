package main

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/Prameesh-P/SHOPRIX/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	initalizers.LoadEnvVariables()
	database.ConnectToDb()
	database.SyncDb()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.New()
	router.Use(gin.Recovery(), middlewares.Logger())
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run()

}
