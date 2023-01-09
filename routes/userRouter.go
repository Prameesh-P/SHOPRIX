package routes

import (
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(ctx *gin.Engine) {
	ctx.POST("/signup", c.Signup)
	ctx.POST("/login", c.Login)
	ctx.GET("/", c.UserHome)
	ctx.POST("/forgetpassword", c.ForgetPassword)
	ctx.GET("/validate", middlewares.RequireAuth, c.Validate)
}
