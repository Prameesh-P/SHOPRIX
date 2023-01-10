package routes

import (
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(ctx *gin.Engine) {

	user := ctx.Group("/user")
	{
		user.GET("/", c.UserHome)

	}
	ctx.POST("/signup", c.Signup)
	ctx.POST("/login", c.Login)
	ctx.POST("/forgetpassword", middlewares.UserAuth(), c.ForgetPassword)
	ctx.GET("/validate", middlewares.UserAuth(), c.Validate)
}
