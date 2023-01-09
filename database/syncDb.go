package database

import "github.com/Prameesh-P/SHOPRIX/models"

func SyncDb() {
	Db.AutoMigrate(&models.User{})
	Db.AutoMigrate(&models.Admin{})
}
