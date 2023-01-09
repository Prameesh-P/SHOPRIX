package database

import "github.com/MohamedmuhsinJ/shopify/models"

func SyncDb() {
	Db.AutoMigrate(&models.User{})
}
