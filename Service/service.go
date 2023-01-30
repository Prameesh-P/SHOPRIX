package Service

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

func SavePayment(charge *models.Charge) (err error) {
	if err = database.Db.Create(charge).Error; err != nil {
		return err
	}
	return nil

}
func Tes(c *gin.Context) {
	c.HTML(200, "app.html", nil)
}
