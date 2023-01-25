package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
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
	userEmail := c.PostForm("user")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	var orderReturn struct {
		OrderId string
	}
	if err := c.ShouldBindJSON(&orderReturn); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	database.Db.Where("orders_id=?", orderReturn.OrderId).Find(&order)
	if order.Order_Status == "returned" {
		c.JSON(400, gin.H{
			"msg": "Item already reaturned",
		})
		c.Abort()
		return
	}
}
