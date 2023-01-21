package database

import "github.com/Prameesh-P/SHOPRIX/models"

func SyncDb() {
	Db.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.ShoeSize{},
		&models.Product{},
		&models.Brand{},
		&models.Category{},
		&models.Address{},
		&models.Cart{},
	)
}
