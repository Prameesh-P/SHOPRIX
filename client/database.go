package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


	var DB *gorm.DB
	var err error

func ConnectToDb() {
	var
	dsn = os.Getenv("DSN")
	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err !=nil{
		log.Fatalf("could't connect databae got an error %v",err)
	}

}
