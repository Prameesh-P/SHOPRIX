package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)
type PageVariable struct {
	OrderId string
	Amount int
	Secret string
}



func RazorPay(c *gin.Context) {
	var user models.User
	var payment models.Charge
	c.BindJSON(&payment)
	useremail := c.Param("user")
	razorPayApiId:=os.Getenv("razorid")
	rezorPaySecret:=os.Getenv("razorsecret")

	database.Db.Raw("select id,phone from users where email=?", useremail).Scan(&user)
	var sumtotal int
	database.Db.Raw("select sum(total_price) from carts where user_id=?", user.ID).Scan(&sumtotal)
	client := razorpay.NewClient(razorPayApiId,rezorPaySecret)
	data := map[string]interface{}{
		"amount":   strconv.Itoa(sumtotal),
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		log.Fatalf("hlo err %v", err)
	}
	value := body["id"]
	str := value.(string)
	homepageVars := PageVariable{
		OrderId: str,
		Amount: sumtotal,
		Secret:rezorPaySecret,
	}

	t, _ := template.ParseFiles("app.html")
	err = t.Execute(c.Writer, homepageVars)
	if err != nil {
		log.Fatalf("err %v", err)
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
