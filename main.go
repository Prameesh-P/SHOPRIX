package main

import (
	"github.com/Prameesh-P/SHOPRIX/routes"
	"os"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run()

}
