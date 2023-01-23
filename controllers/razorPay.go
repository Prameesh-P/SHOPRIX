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
	userEmail:=c.GetString("user")
	database.Db.Raw("select id,phone from users where email=?",userEmail).Scan(&user)
	var sumtotal uint
	database.Db.Raw("select sum(total_prize) from carts where user_id=?",user.ID).Scan(&sumtotal)
	fmt.Println(sumtotal)

	client := razorpay.NewClient("rzp_test_7eQ9M4ASOPr3Ul", "k4vBqTSaaOWHUYaikHwCR3S7")	
	razorPayValue:=sumtotal*100
	data := map[string]interface{}{
		"amount":   razorPayValue,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	
}