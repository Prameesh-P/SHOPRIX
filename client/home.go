package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Prameesh-P/SHOPRIX/messages"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserHome(c *gin.Context) {
	cc,err:=grpc.Dial("localhost"+port,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if  err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Err": "Error whiler dailing",
		})
		return
	}
	defer cc.Close()

	client := messages.NewAuthentifiationServiceClient(cc)
	req:=&messages.HomeRequest{}
	res,err:=client.UserHomeService(context.Background(),req)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"Err":"Did't get the response from service please check your connection",
		})
		return
	}

	fmt.Println(res.Res)
	c.JSON(http.StatusOK,gin.H{
		"Success":res.Res,
	})

}