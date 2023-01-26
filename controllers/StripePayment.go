package controllers

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"net/http"
	"os"

	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

func Stripe(c *gin.Context) {
	var json models.Charge
	c.BindJSON(&json)

	// Set Stripe API key
	apiKey := os.Getenv("STRIPE_KEY")
	stripe.Key = apiKey
	_, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(json.Amount),
		Currency:     stripe.String(string(stripe.CurrencyINR)),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")}, // this should come from clientside
		ReceiptEmail: stripe.String(json.ReceiptEmail)})

	if err != nil {
		// Handle any errors from attempt to charge
		c.String(http.StatusBadRequest, "Request failed")
		return
	}
	c.String(http.StatusCreated, "Successfully charged")
}
