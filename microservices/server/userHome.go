package main

import (
	"context"
	"fmt"

	"github.com/Prameesh-P/SHOPRIX/messages"
)

func (s *Server) UserHomeService(ctx context.Context,req *messages.HomeRequest)(*messages.HomeResponse,error){
	
	fmt.Println("User on Home Page..")
	return &messages.HomeResponse{
		Res: "Welcome to User Home Page",
	},nil
	
}