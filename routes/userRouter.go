package routes

import (
	c "github.com/Prameesh-P/E-COMMERCE/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(ctx *gin.Engine) {
	ctx.POST("/signup", c.SignUp)
	ctx.POST("/login", c.Login)
	ctx.POST("/forgetpassword", c.ForgetPassword)
}
