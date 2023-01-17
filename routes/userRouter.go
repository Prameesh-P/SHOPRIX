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
		user.GET("/profile",c.UserProfileGet)
		user.POST("/profile/edit",c.EditUserAddress)
		user.POST("/profile/add",c.AddAddress)
	}
	routes.POST("/signup", c.Signup)
	routes.POST("/login", c.Login)
	routes.GET("/forgetpassword", middlewares.UserAuth(), c.ForgetPassword)
	routes.GET("/validate", middlewares.UserAuth(), c.Validate)
	routes.POST("/login/otp", c.OtpLog)
	routes.POST("/login/otpvalidate", c.CheckOTP)
	routes.GET("/logout", middlewares.AdminAuth())
	routes.GET("/forgetemail/:email",middlewares.UserAuth(),c.ForgetPasswordEmail)
}
