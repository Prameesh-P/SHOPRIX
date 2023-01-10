package middlewares

import (
	"github.com/Prameesh-P/SHOPRIX/authentification"
	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		// tokenString := context.GetHeader("Authorization")
		tokenString, err := context.Cookie("Adminjwt")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err = authentification.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
func UserAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, err := context.Cookie("UserAuth")
		if tokenString == "" {
			context.JSON(401, gin.H{
				"error": "Request does not contain an access token",
			})
			context.Abort()
			return
		}
		err = authentification.ValidateToken(tokenString)
		context.Set("user", authentification.P)
		if err != nil {
			context.JSON(401, gin.H{
				"error": err.Error(),
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
