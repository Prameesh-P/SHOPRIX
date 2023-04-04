package main

import(

	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine){

	c.GET("/",UserHome)
	c.POST("/login",Login)
	c.POST("/signup",Signup)

}