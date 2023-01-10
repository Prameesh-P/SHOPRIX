package routes

import (
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(ctx *gin.Engine) {
	admin := ctx.Group("/admin")
	{
		admin.POST("/signup", c.AdminSignup)
		admin.POST("/login", c.AdminLogin)
		admin.GET("/", middlewares.AdminAuth(), c.AdminHome)
		admin.GET("/userdata", middlewares.AdminAuth(), c.Userdata)
		admin.PUT("/userdata/block/:id", middlewares.AdminAuth(), c.BlockUser)
	}
}
