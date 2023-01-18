package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var user models.User
	var product models.Product
	userEmail:=c.GetString("user")
	database.Db.Raw("select id from users where email=?",userEmail).Scan(&user)
	var ProductDetails struct{
		ProductID uint
		Quantity  uint
	}
	c.ShouldBindJSON(&ProductDetails)
	database.Db.Raw("select price,stock from products where product_id=?",ProductDetails.ProductID).Scan(&product)
	total:=product.Price*ProductDetails.Quantity
	productId:=ProductDetails.ProductID
	productQuantity:=ProductDetails.Quantity
	cart:=models.Cart{
		ProductID: productId,
		Quantity: productQuantity,
		UserId: user.ID,
		TotalPrice: total,
	}
	var Cart []models.Cart 
}