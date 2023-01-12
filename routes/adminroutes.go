package routes

import (
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.Engine) {
	admin := routes.Group("/admin")

	admin.POST("/signup", c.AdminSignup)
	admin.POST("/login", c.AdminLogin)
	admin.GET("/", middlewares.AdminAuth(), c.AdminHome)
	admin.GET("/userdata", middlewares.AdminAuth(), c.UserData)
	admin.PUT("/userdata/block/:id", middlewares.AdminAuth(), c.BlockUser)
	admin.PUT("/userdata/unblock/:id", middlewares.AdminAuth(), c.UnBlockUser)
	admin.GET("/getcategory", middlewares.AdminAuth(), c.ListAllCategory)
	admin.POST("/addproducts", middlewares.AdminAuth(), c.ProductAdding)
	admin.POST("/editproducts", middlewares.AdminAuth(), c.EditProducts)
}
