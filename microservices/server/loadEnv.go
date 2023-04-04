package main

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	if err := godotenv.Load();err!=nil{
		log.Fatalf("Could't load env variable..Got an err :%v",err)	
	}

}