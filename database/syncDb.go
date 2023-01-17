package database

import "github.com/Prameesh-P/SHOPRIX/models"

func SyncDb() {
	Db.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.ShoeSize{},
		&models.Product{},
		&models.Brand{},
		&models.WishList{},
		&models.Category{},
		&models.Address{},
	)
}
