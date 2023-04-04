package main

import (
	"context"
	"log"

	msg "github.com/Prameesh-P/SHOPRIX/messages"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) SignUpService(ctx context.Context, req *msg.SignupRequest) (*msg.SignupRespone, error) {
	firstName := req.GetFirstName()
	lastName := req.GetLastName()
	email := req.GetEmail()
	password := req.GetPassword()
	phone := req.GetPhone()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("error got while password hashing %v", err)
	}
	users := &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(bytes),
		Phone:     phone,
	}
	record := DB.Create(&users)
	if record.Error != nil {
		log.Fatalf("error could not create database%v", err)
	}
	log.Printf("User %s Signuped with %s", firstName, email)
	return &msg.SignupRespone{
		Name:   firstName,
		Email:  email,
		Result: "success",
	}, nil
}
