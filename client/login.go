package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Prameesh-P/SHOPRIX/database"
	msg "github.com/Prameesh-P/SHOPRIX/messages"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Login(c *gin.Context) {
	var Body struct {
		Email    string
		Password string
	}
	if c.ShouldBind(&Body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get request",
		})
		return
	}
	var user models.User
	database.Db.First(&user, "email=?", Body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}
	if user.BlockStatus {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "user has been Blocked by admin",
		})
		c.Abort()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("Secret")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserJWT", tokenString, 3600*24*30, "", "", false, true)

	fmt.Println(Body.Password)
	cc, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Err": "Error while dailing",
		})
		return
	}
	defer cc.Close()
	client := msg.NewAuthentifiationServiceClient(cc)
	req := &msg.LoginRequest{
		Email:    Body.Email,
		Password: Body.Password,
	}
	res, err := client.LoginService(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "Something went wrong your request please check your given details is correct or not !!!",
		})
		c.Abort()
		return
	}
	fmt.Println(res)
	c.JSON(http.StatusOK, gin.H{
		"success": "OK",
		"token":   tokenString,
	})
}
