package controllers

import (
	"math/rand"
	"strconv"
	"time"
)

type Orderd_Items []struct {
	UserId          uint
	Product_id      uint
	OrdersID        string
	Product_Name    string
	Price           string
	Order_Status    string
	Payment_Status  string
	PaymentMethod   string
	Applied_Coupons string
	Total_amount    uint
}

func CreateOrderId() string {

	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id
	return orderID
}
