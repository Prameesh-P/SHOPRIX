package routes

import (
	"github.com/Prameesh-P/SHOPRIX/Service"
	// "github.com/Prameesh-P/SHOPRIX/controllers"
	c "github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/middlewares"
	"github.com/gin-gonic/gin"
	
)


func UserRoutes(routes *gin.Engine) {
	user := routes.Group("/user")


	routes.POST("/signup",c.Signup)
	routes.POST("/login",c.Login)
	routes.GET("/",c.UserHome)
	// routes.GET("/otp",controllers.NexmoOtpVErification)
	
	/*------------------------Group Routes >> user--------------------------- */


	{
		user.GET("/profile", middlewares.UserAuth(), c.UserProfileGet)
		user.PUT("/profile/edit", middlewares.UserAuth(), c.EditUserAddress)
		user.POST("/profile/add", middlewares.UserAuth(), c.AddAddress)
		user.GET("/show-product-id", middlewares.UserAuth(), c.ShowProductsID)
		user.GET("/get-productbyid", middlewares.UserAuth(), c.GetProductByID)
		user.GET("/view-products", middlewares.UserAuth(), c.ProductView)
		user.POST("/addtocart", middlewares.UserAuth(), c.AddToCart)
		user.GET("/viewcart", middlewares.UserAuth(), c.ViewCart)
		user.GET("/payment/:user", middlewares.UserAuth(),c.RazorPay)
		user.POST("/checkoutAddress", middlewares.UserAuth(), c.CheckOutAddress)
		user.GET("/checkout", middlewares.UserAuth(), c.CheckOut)
		user.GET("/vieworder", middlewares.UserAuth(), c.ViewOrders)
		user.GET("/returnorder", middlewares.UserAuth(), c.ReturnOrders)
		user.GET("/cancelorder", middlewares.UserAuth(), c.CancelOrders)
	}


	/*---------------------------without Groups---------------------------------*/



	routes.GET("/forgetpassword", middlewares.UserAuth(), c.ForgetPassword)
	routes.GET("/validate", middlewares.UserAuth(), c.Validate)
	routes.POST("/login/otp", c.OtpLog)
	routes.POST("/login/otpvalidate", c.CheckOTP)
	routes.GET("/logout", middlewares.AdminAuth())
	routes.GET("/forgetemail/:email", middlewares.UserAuth(), c.ForgetPasswordEmail)
	routes.GET("/app", Service.Tes)
	routes.POST("/",c.Signup)

	
	
	/*---------------------MicorService Routes----------------------*/


}
