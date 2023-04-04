package main

import (
	"context"
	"log"
	msg"github.com/Prameesh-P/SHOPRIX/messages"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) LoginService(ctx context.Context, req *msg.LoginRequest) (*msg.LoginResponse, error){
	email:=req.GetEmail()
	password:=req.GetPassword()
	var k User
    DB.Where("email = ?", email).First(&k)
	err:=bcrypt.CompareHashAndPassword([]byte(k.Password),[]byte(password))
	if err!=nil{
		log.Printf("user password does not match %v",err)
		return nil,err
	}

log.Printf("User %s login with %s",k.FirstName,email)
return &msg.LoginResponse{
Name: k.FirstName,
Email: email,
Result: "success",
},nil

}