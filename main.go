package main

import (
	"os"

	"github.com/MohamedmuhsinJ/shopify/controllers"
	"github.com/MohamedmuhsinJ/shopify/database"
	"github.com/MohamedmuhsinJ/shopify/initalizers"
	"github.com/MohamedmuhsinJ/shopify/middlewares"
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
	// controllers.Otp("9159564424")
	// controllers.CheckOtp("9159564424")
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middlewares.RequireAuth, controllers.Validate)
	router.Run()

	// routes.AuthRoutes(router)
	// routes.UserRoutes(router)
}
