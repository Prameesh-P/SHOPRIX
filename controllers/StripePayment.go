package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"net/http"
	"os"
)

func Stripe(c *gin.Context) {
	var user models.User
	var payment models.Charge
	c.BindJSON(&payment)
	useremail := c.Request.FormValue("user")

	// Set Stripe API key
	apiKey := os.Getenv("STRIKE_KEY")
	stripe.Key = apiKey
	database.Db.Raw("select id,phone from users where email=?", useremail).Scan(&user)
	var sumtotal int
	database.Db.Raw("select sum(total_price) from carts where user_id=?", user.ID).Scan(&sumtotal)

	_, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(payment.Amount),
		Currency:     stripe.String(string(stripe.CurrencyINR)),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: stripe.String("prameepramee0@gmail.com")})

	if err != nil {
		c.String(http.StatusBadRequest, "Request failed")
		return
	}
	c.HTML(200, "app.html", gin.H{

		"UserID":       user.ID,
		"total_price":  sumtotal,
		"total":        sumtotal,
		"amount":       sumtotal,
		"Email":        useremail,
		"Phone_Number": user.Phone,
	})
	c.String(http.StatusCreated, "Successfully charged")
}

//var user models.User
//// var cart models.Cart
//useremail := c.Request.FormValue("user")
//fmt.Println(useremail)
//database.Db.Raw("select id,phone from users where email=?", useremail).Scan(&user)
//var sumtotal uint
//database.Db.Raw("select sum(total_price) from carts where user_id=?", user.ID).Scan(&sumtotal)
//
//fmt.Println(sumtotal)
//client := razorpay.NewClient("rzp_test_7eQ9M4ASOPr3Ul", "k4vBqTSaaOWHUYaikHwCR3S7")
//razpayvalue := sumtotal * 100
//data := map[string]interface{}{
//	"amount":   razpayvalue,
//	"currency": "INR",
//	"receipt":  "some_receipt_id",
//}
//body, err := client.Order.Create(data, nil)
//fmt.Println(body)
//value := body["id"]
//
//if err != nil {
//	c.JSON(404, gin.H{
//		"msg": "error creating order",
//	})
//}
//c.HTML(200, "app.html", gin.H{
//
//	"UserID":       user.ID,
//	"total_price":  sumtotal,
//	"total":        razpayvalue,
//	"orderid":      value,
//	"amount":       sumtotal,
//	"Email":        useremail,
//	"Phone_Number": user.Phone,
//})
//if err != nil {
//	c.JSON(200, gin.H{
//		"msg": value,
//	})
//}
//}
//func RazorpaySuccess(c *gin.Context) {
//	userid := c.Query("user_id")
//	userID, _ := strconv.Atoi(userid)
//	orderid := c.Query("order_id")
//	paymentid := c.Query("payment_id")
//	signature := c.Query("signature")
//	id := c.Query("orderid")
//	totalamount := c.Query("total")
//	Rpay := models.RazorPay{
//		UserID:          uint(userID),
//		OrderId:         id,
//		RazorPaymentId:  paymentid,
//		Signature:       signature,
//		RazorPayOrderID: orderid,
//		AmountPaid:      totalamount,
//	}
//	initializers.DB.Create(&Rpay)
//	var cart models.Cart
//	initializers.DB.Raw("delete from carts where user_id=?", userID).Scan(&cart)
//	fmt.Println(userID, orderid)
//	OrderPlaced(userID, orderid)
//
//	c.JSON(200, gin.H{
//
//		"status": true,
//	})
//
//}
//func Success(c *gin.Context) {
//	c.HTML(200, "succs.html", nil)
//
//}
//func OrderPlaced(Uid int, orderId string) {
//	userid := Uid
//	orderid := orderId
//	var orders models.Orders
//	var applied string
//
//	initializers.DB.Raw("update orders set order_status=?,payment_status=?,order_id=? where user_id=?", "order completed", "payment done", orderid, userid).Scan(&orders)
//	var ordereditems models.Orderd_Items
//	initializers.DB.Raw("select applied_coupons from orders where order_id=?", orderId).Scan(&applied)
//	coupons := models.Applied_Coupons{
//		Coupon_Code: applied,
//		UserID:      uint(userid),
//	}
//	initializers.DB.Create(&coupons)
//	initializers.DB.Raw("update orderd_items set order_status=?,payment_status=? where user_id=?", "orderplaced", "Payment Completed", userid).Scan(&ordereditems)
//
//}
