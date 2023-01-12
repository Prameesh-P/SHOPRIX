package database

import "github.com/Prameesh-P/SHOPRIX/models"

func SyncDb() {
	Db.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Otp{},
		&models.ShoeSize{},
		&models.Product{},
		&models.Brand{},
		&models.WishList{},
		&models.Category{},
	)
}
