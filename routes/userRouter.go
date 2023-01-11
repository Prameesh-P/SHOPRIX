package routes

import (
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {

	user := routes.Group("/user")
	{
		user.GET("/", c.UserHome)
	}
	routes.POST("/signup", c.Signup)
	routes.POST("/login", c.Login)
	routes.POST("/forgetpassword", middlewares.UserAuth(), c.ForgetPassword)
	routes.GET("/validate", middlewares.UserAuth(), c.Validate)
	routes.POST("/login/otp", c.OtpLog)
	routes.POST("/login/otpvalidate", c.CheckOTP)
}
