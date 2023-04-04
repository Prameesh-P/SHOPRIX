package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Prameesh-P/SHOPRIX/database"
	msg "github.com/Prameesh-P/SHOPRIX/messages"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":6000"
)

func Signup(c *gin.Context) {
	var Body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Phone     string `json:"phone"`
	}
	// var user models.User
	if c.ShouldBindJSON(&Body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Please enter input correctly",
		})
		return
	}
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "coult not connect server",
		})
		return
	}
	defer conn.Close()

	cli := msg.NewAuthentifiationServiceClient(conn)
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		c.Abort()
		return
	}
	user := models.User{FirstName: Body.FirstName, LastName: Body.LastName, Email: Body.Email, Password: string(hash), Phone: Body.Phone}
	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		c.Abort()
		return
	}

	req := &msg.SignupRequest{
		FirstName: Body.FirstName,
		LastName:  Body.LastName,
		Email:     Body.Email,
		Password:  Body.Password,
		Phone: Body.Phone,
	}
	res, err := cli.SignUpService(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Error from response",
		})
		return
	}

	var resp struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		result string `json:"result"`
	}
	fmt.Println(res)
	resp.result = res.Result
	resp.Email = res.Email
	resp.Name = res.Name
	c.JSON(http.StatusOK, gin.H{
		"success": "user signup Successfully completed..",
	})

}
