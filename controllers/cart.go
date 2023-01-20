package controllers

import (
	"fmt"
	"net/http"

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
	database.Db.Raw("select cart_id,product_id from carts where user_id=?",user.ID).Scan(&Cart)
	for _,v:=range Cart{
		fmt.Println("loop started..!")
		if v.ProductID==productId {
			fmt.Println("in the condition..")
			database.Db.Raw("select quantity from carts where product_id=? and user_id=?",ProductDetails.ProductID,user.ID).Scan(&Cart)
			totlV:=(productQuantity+cart.Quantity)*product.Price
			database.Db.Raw("update carts set quantity=?,total_price where  product_id=? and user_id=?",productQuantity+cart.Quantity,totlV,productId,user.ID)
			c.JSON(http.StatusOK,gin.H{
				"msg":"quantiity updated successfully",
			})
			c.Abort()
			return
		}
	}
	
}