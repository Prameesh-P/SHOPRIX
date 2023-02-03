package middlewares

import (
	"github.com/Prameesh-P/SHOPRIX/authentification"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {

		tokenString, err := context.Cookie("AdminJWT")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err = authentification.ValidateToken(tokenString)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
func UserAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, err := context.Cookie("UserJWT")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Request does not contain an access token",
			})
			context.Abort()
			return
		}
		err = authentification.ValidateToken(tokenString)
		context.Set("user", authentification.P)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
