package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)


// @Summary admin-coupen
// @ID admin-userdata-coupen
// @Description admin userdata
// @Tags Admin
// @Accept json
// @Produce json
// @Param coupen_code formData string true "coupen code"
// @Param discount formData string true "discount of the coupen"
// @Param quantity formData string true "quantity code"
// @Param validity formData string true "validity code"
// @Success 200 
// @Failure 400 
// @Router /admin/generate-coupens/ [post]
func GenerateCoupens(c *gin.Context) {
	coupenCode := c.PostForm("coupen_code")
	coupenDiscount := c.PostForm("discount")
	discount, _ := strconv.Atoi(coupenDiscount)
	coupenQuantity := c.PostForm("quantity")
	quantity, _ := strconv.Atoi(coupenQuantity)
	coupenValidity := c.PostForm("validity")
	validity, _ := strconv.Atoi(coupenValidity)

	expirationTime := time.Now().AddDate(0, 0, validity)
	expirationTime.Unix()

	coupens := models.Coupon{
		Coupon_Code: coupenCode,
		Discount:    uint(discount),
		Quantity:    uint(quantity),
		Validity:    int64(validity),
	}

	database.Db.Create(&coupens)

	c.JSON(http.StatusOK, gin.H{
		"coupon":      coupens.Coupon_Code,
		"coupon-code": coupens.Coupon_Code,
		"msg":         "coupon added",
	})
}
