package controllers

import (
	"fmt"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
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
func ViewOrders(c *gin.Context) {
	var user models.User
	var ordered_items Orderd_Items
	userEmail := c.PostForm("user")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	record := database.Db.Raw("select user_id,product_id,product_name,applied_coupons,price,orders_id,order_status,payment_status,payment_method,total_amount from ordered_items where user_id=?", user.ID).Scan(&ordered_items)
	if search := c.PostForm("search"); search != "" {
		query := database.Db.Raw("select user_id,product_id,product_name,applied_coupons,price,orders_id,order_status,payment_status,payment_method,total_amount from ordered_items where (product_name ilike ? or payment_method ilike ? )and user_id=? ", "%"+search+"%", "%"+search+"%", user.ID).Scan(&ordered_items)
		if query.Error != nil {
			c.JSON(404, gin.H{
				"err": query.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"orders": ordered_items,
	})
}
func ReturnOrders(c *gin.Context) {
	var order models.OrderedItems
	var user models.User
	userEmail := c.Query("user")
	fmt.Println(userEmail)
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	var orderReturn struct {
		OrderId string
	}
	if err := c.BindJSON(&orderReturn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	database.Db.Where("orders_id=?", orderReturn.OrderId).Find(&order)
	if order.Order_Status == "returned" {
		c.JSON(400, gin.H{
			"msg": "Item already returned",
		})
		c.Abort()
		return
	}
	var balance int
	database.Db.Raw("select wallet_balance from wallets where user_id=?", user.ID).Scan(&balance)
	//tx := database.Db.Begin()
	record := database.Db.Model(&models.OrderedItems{}).Where("orders_id=?", orderReturn.OrderId).Update("order_status", "returned")
	if record.Error != nil {
		//tx.Rollback()
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	//query := database.Db.Raw("delete from ordered_items where order_id=?",orderReturn.OrderId)
	//if query.Error != nil {
	//	c.JSON(400, gin.H{
	//		"error": query.Error.Error(),
	//	})
	//	c.Abort()
	//	return
	//}
	newBalance := balance + int(order.Total_amount)
	record1 := database.Db.Model(&models.Wallet{}).Where("user_id=?", user.ID).Update("wallet_balance", newBalance)
	if record1.Error != nil {
		//tx.Rollback()
		c.JSON(404, gin.H{
			"err": record1.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"msg": "order returned",
	})
}
func CancelOrders(c *gin.Context) {
	var user models.User
	var returned string
	var returns = "returned"
	userEmail := c.Query("user")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	oderID := c.Query("orderid")
	var orders models.OrderedItems
	database.Db.Where("orders_id=?", oderID).Find(&orders)
	fmt.Println(orders.Order_Status, orders.OrdersID)
	fmt.Println(oderID)
	database.Db.Raw("select order_status from ordered_items where order_status=?", returns).Scan(&returned)
	if returned != "" {
		c.JSON(400, gin.H{
			"message": "Can't cancel products..!! Because order is already returned",
		})
		c.Abort()
		return
	}
	if orders.Order_Status == "order cancelled" {
		c.JSON(400, gin.H{
			"status": "false",
			"msg":    "Order already cancelled..!!",
		})
		return
	}
	database.Db.Raw("update ordered_items set order_status=? where orders_id=?", "order cancelled", oderID).Scan(&orders)
	var price int
	database.Db.Raw("SELECT total_amount FROM ordered_items WHERE orders_id = ?", oderID).Scan(&price)
	var balance int
	database.Db.Raw("SELECT wallet_balance FROM wallets WHERE id = ?", user.ID).Scan(&balance)
	newBalance := price + balance
	database.Db.Raw("update wallets set wallet_balance=? where id=?", newBalance, user.ID).Scan(&user)
	c.JSON(200, gin.H{
		"msg": "order cancelled",
	})
}
