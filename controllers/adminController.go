package controllers

import "github.com/gin-gonic/gin"

func AdminSignup(c *gin.Context) {
	var Admin struct{
		Email string
		Password string
	}
	if err:=c.ShouldBind(&Admin);err!=nil {
		c.JSON(404,gin.H{
			"error":err.Error(),
		})
	}
	
}