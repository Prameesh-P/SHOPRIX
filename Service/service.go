package Service

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
)

func SavePayment(charge *models.Charge) (err error) {
	if err = database.Db.Create(&charge).Error; err != nil {
		return err
	}
	return nil

}
