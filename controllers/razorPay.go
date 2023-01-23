package controllers

import (
	"fmt"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	var user models.User
	userEmail := c.GetString("user")
	database.Db.Raw("select id,phone from users where email=?", userEmail).Scan(&user)
	var sumtotal uint
	database.Db.Raw("select sum(total_prize) from carts where user_id=?", user.ID).Scan(&sumtotal)
	fmt.Println(sumtotal)

	client := razorpay.NewClient("rzp_test_7eQ9M4ASOPr3Ul", "k4vBqTSaaOWHUYaikHwCR3S7")
	razorPayValue := sumtotal * 100
	data := map[string]interface{}{
		"amount":   razorPayValue,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	fmt.Println(body)
	value := body["id"]

	if err != nil {
		c.JSON(404, gin.H{
			"msg": "error creating order",
		})
	}
	c.HTML(200, "app.html", gin.H{

		"UserID":       user.ID,
		"total_price":  sumtotal,
		"total":        razorPayValue,
		"orderid":      value,
		"amount":       sumtotal,
		"Email":        userEmail,
		"Phone_Number": user.Phone,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"msg": value,
		})
	}
}
