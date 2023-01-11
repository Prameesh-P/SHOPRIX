package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func ConnectToDb() {
	var err error
	DBHost := os.Getenv("dbHost")
	DBUser := os.Getenv("dbUser")
	DBPassword := os.Getenv("dbPassword")
	DBName := os.Getenv("dbName")
	DDPort := os.Getenv("dbPort")
	dsn := fmt.Sprintf("host=%s user=%s password=%s DBName=%s sslmode=disable port=%s", DBHost, DBUser, DBPassword, DBName, DDPort)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db not connected")
	}

}
