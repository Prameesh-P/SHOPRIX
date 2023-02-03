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
	c.HTML(http.StatusOK, "app.html", gin.H{

		"UserID":       user.ID,
		"total_price":  sumtotal,
		"total":        sumtotal,
		"amount":       sumtotal,
		"Email":        useremail,
		"Phone_Number": user.Phone,
	})
	c.String(http.StatusCreated, "Successfully charged")
}
